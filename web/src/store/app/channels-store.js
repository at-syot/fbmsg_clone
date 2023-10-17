import { writable } from "svelte/store";
import _ from "lodash/fp";
import dayjs from "dayjs";
import { get } from "svelte/store";
import * as persistedKeys from '../persisted-keys.js'
import { websocketMessageStore, websocketStore } from "./websocket-store";
import { uiSidebarDisplayModeStore } from "../ui/sidebar-display-store";
import { uiContentPanelDisplayStore } from "../ui/content-pannel-display-store";
import { serverHost } from '../../lib/client'

const PERSISTED_KEY = persistedKeys.USER_CHANNELS_KEY;

function createChannelsStore() {
  const { subscribe, set, update } = writable([]);

  return {
    subscribe,
    getPersisted() {
      const persisted = localStorage.getItem(PERSISTED_KEY);
      if (persisted) {
        const parsed = JSON.parse(persisted);
        set(parsed);
        return parsed;
      }
    },
    persist() {
      update((data) => {
        // remove channel's messages
        const withoutChanMsgs = _.pipe(
          _.cloneDeep,
          _.map((ch) => {
            if (!ch.messages) return ch;
            ch.messages = [];
            return ch;
          })
        )(data);

        localStorage.setItem(PERSISTED_KEY, JSON.stringify(withoutChanMsgs));
        return data;
      });
    },

    /** 
    * @function getActiveChannel
    * @returns {Object} 
    */
    getActiveChannel(channels) {
      const activeChans = channels.filter((ch) => ch.active);
      if (activeChans.length === 0) return;
      return activeChans[0];
    },

    /** @function setActiveChannel */
    async setActiveChannel(channelId) {
      const messages = await fetchChannelMessages(channelId);
      update(
        _.pipe(
          _.map((ch) => {
            return _.set("active", ch.id == channelId)(ch);
          }),
          _.map((ch) => {
            if (ch.id !== channelId) return ch;

            const msgItem = (() => {
              const _msgItem = messages[messages.length - 1];
              if (_.isEmpty(_msgItem))
                return { id: "", message: "-", createdAt: "-" };

              const { createdAt } = _msgItem;
              return { ..._msgItem, createdAt: dayjs(createdAt).format("ddd") };
            })();
            ch.latestMsgItem = msgItem;
            ch.messages = messages;

            return ch;
          })
        )
      );

      this.persist();
    },

    /**
     * @function getChannelByUsers
     * CHEATING here, move this login to backend
     *
     * result mush have only one channel, if not: have duplicated channels
     */
    getChannelByUsers(userIds, channels) {
      const existingChans = channels.filter(
        (ch) =>
          ch.users.every((u) => userIds.indexOf(u.id) !== -1) &&
          ch.users.length === userIds.length
      );
      if (existingChans.length > 1) {
        throw Error("found duplicated channels");
      }
      return existingChans[0];
    },

    /**
     * @function createAndAddNewChannel
     *  - create new channel
     *  - fetch channel info item
     *  - add to store
     */
    async createAndAddNewChannel(creatorId, userIds, displayname) {
      const existingChan = this.getChannelByUsers(
        [creatorId, ...userIds],
        get(channelsStore)
      );
      let chanItem = existingChan;

      if (!existingChan) {
        const { channelId } = await createChannel(
          creatorId,
          userIds,
          displayname
        );
        const _chanItem = await fetchChannelById(channelId);
        const chanItems = [withDefaultProps(creatorId)(_chanItem)];

        update(_.concat(chanItems));
        chanItem = _chanItem;
      }

      const { id: chanId } = chanItem;
      await websocketStore.joinChannel(chanId, creatorId);
      uiSidebarDisplayModeStore.setSidebarDisplayMode("channels_list");
      websocketMessageStore.clearMessage();
      await this.setActiveChannel(chanId);
      uiContentPanelDisplayStore.setDisplaymode("message");

      return chanItem;
    },

    pushChannelMessage(channelId, messageItem) {
      if (!channelId || _.isEmpty(messageItem)) return;
      update(
        _.map((ch) => {
          if (ch.id !== channelId) return ch;
          ch.latestMsgItem = {
            ...messageItem,
            createdAt: dayjs(messageItem.createdAt).format("ddd"),
          };
          ch.messages.push(messageItem);
          return ch;
        })
      );

      this.persist();
    },

    /**
     * @function fetchChannels
     * - fetch user channels
     * - set default props
     * - merge with persisted channels if need
     */
    async fetchChannels(userId) {
      const channels = await fetchChannels(userId);
      const mappedDefaultChannels = _.map(withDefaultProps(userId))(channels);

      let toSetChans = [];
      const persistedChans = JSON.parse(localStorage.getItem(PERSISTED_KEY));
      if (persistedChans) {
        const mergedWithPersistedChs = mappedDefaultChannels.map((ch) => {
          const persistedCh = persistedChans.filter(
            ({ id }) => id === ch.id
          )[0];
          return _.merge(ch, persistedCh);
        });
        toSetChans = mergedWithPersistedChs;
      } else {
        toSetChans = mappedDefaultChannels;
      }

      set(toSetChans);
      this.persist();

      return toSetChans;
    },

    /** 
    * @function autoJoinChannel
    * @param {string} userId
    */
    async autoJoinChannel(userId) {
      const channels = await this.fetchChannels(userId)
      if (!channels || channels.length == 0) return

      // auto join first channel
      const { id: channelId } = channels[0]
      await this.setActiveChannel(channelId)
      await websocketStore.joinChannel(channelId, userId)
    },

    clearChannels() {
      set([])
    }
  };
}

function withDefaultProps(userId) {
  const displayname = (channel) => {
    const { displayname } = channel;
    if (displayname) return displayname;

    // users.length === 2 --> private chat
    const { users } = channel;
    if (users.length === 2) {
      const { username } = users.filter((u) => u.id !== userId)[0];
      return username;
    }

    return users.map((u) => u.username).join(" <> ");
  };

  return (channel) => {
    const setDisplayName = _.set("displayname", displayname(channel));
    const setUnActive = _.set("active", false);
    const setEmptyMessages = _.set("messages", []);
    const setLatestMessage = _.set("latestMsgItem", {
      message: "",
      createdAt: "",
    });
    const setDefaultChannelProps = _.pipe(
      setDisplayName,
      setUnActive,
      setLatestMessage,
      setEmptyMessages
    );
    return setDefaultChannelProps(channel);
  };
}

async function createChannel(creatorId, userIds, displayname) {
  const endpoint = `${serverHost()}/channels`;
  const body = JSON.stringify({ creatorId, userIds, displayname });
  const res = await fetch(endpoint, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body,
  });
  if (!res.ok) {
    alert(`creating channel err - ${await res.text()}`);
    return;
  }
  return await res.json();
}

async function fetchChannels(userId) {
  const url = `${serverHost()}/users/${userId}/channels`;
  const res = await fetch(url);
  if (!res.ok) {
    console.log("err -", await res.text());
    return;
  }
  const { channels } = await res.json();
  return channels;
}

async function fetchChannelById(channelId) {
  const res = await fetch(`${serverHost()}/channels/${channelId}`);
  if (!res.ok) {
    console.log("fetch channel by id err - ", await res.text());
    return;
  }
  return await res.json();
}

async function fetchChannelMessages(channelId) {
  const res = await fetch(
    `${serverHost()}/channels/${channelId}/messages`
  );
  if (!res.ok) {
    console.log(`fetch channel messages err - ${await res.text()}`);
    return;
  }
  const { messages } = await res.json();
  return messages;
}

export const channelsStore = createChannelsStore();
