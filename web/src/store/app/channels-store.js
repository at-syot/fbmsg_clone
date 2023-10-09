import {writable} from "svelte/store";
import _ from "lodash/fp";
import dayjs from "dayjs";
import {userStore} from "./user-store.js";

const PERSISTED_KEY = "app-user-channels";

function createChannelsStore() {
  const {subscribe, set, update} = writable([]);
  
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
        localStorage.setItem(PERSISTED_KEY, JSON.stringify(data));
        return data;
      });
    },
    async setActiveChannel(channelId) {
      const messages = await this.fetchChannelMessages(channelId)
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
                return {id: '', message: "-", createdAt: "-"};
              
              const {createdAt} = _msgItem;
              return {..._msgItem, createdAt: dayjs(createdAt).format("ddd")}
            })();
            
            return _.pipe(
              _.set("latestMsgItem", msgItem),
              _.set("messages", messages)
            )(ch);
          })
        )
      );
      
      this.persist()
    },
    
    /** - create new channel
     *  - fetch channel info item
     *  - add to store
     *  TODO: set default channelItem's values
     */
    async createAndAddNewChannel(creatorId, userIds) {
      const {channelId} = await createChannel(creatorId, userIds)
      const chanItem = await fetchChannelById(channelId)
      const chanItems = [withDefaultProps(creatorId)(chanItem)]
      update(_.concat(chanItems))
    },
    pushChannelMessage(channelId, messageItem) {
      console.log('push ', channelId, messageItem)
      if (!channelId || _.isEmpty(messageItem)) return;
      update(
        _.map((ch) => {
          if (ch.id !== channelId) return ch;
          ch.latestMsgItem = {...messageItem, createdAt: dayjs(messageItem.createdAt).format("ddd")};
          ch.messages.push(messageItem);
          return ch;
        })
      );
      this.persist();
    },
    async fetchChannels(userId) {
      const channels = await fetchChannels(userId)
      console.log('channels', channels)
      const mappedDefaultChannels = _.map(withDefaultProps(userId))(channels);
      localStorage.setItem(
        PERSISTED_KEY,
        JSON.stringify(mappedDefaultChannels)
      );
      set(mappedDefaultChannels);
    },
    async fetchChannelMessages(channelId) {
      const res = await fetch(
        `http://localhost:3000/channels/${channelId}/messages`
      );
      if (!res.ok) {
        console.log(`fetch channel messages err - ${await res.text()}`);
        return;
      }
      const {messages} = await res.json();
      return messages;
    },
  };
}

function withDefaultProps(userId) {
  const displayname = (channel) => {
    const {displayname} = channel
    if (displayname) return displayname
    
    // users.length === 2 --> private chat
    const {users} = channel
    if (users.length === 2) {
      const {username} = users.filter(u => u.id !== userId)[0]
      return username
    }
    
    return users.map(u => u.username).join(' <> ')
  }
  
  return (channel) => {
    const setDisplayName = _.set('displayname', displayname(channel))
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
    return setDefaultChannelProps(channel)
  }
}

async function createChannel(creatorId, userIds) {
  const endpoint = "http://localhost:3000/channels"
  const body = JSON.stringify({creatorId, userIds})
  const res = await fetch(endpoint, {
    method: 'POST',
    headers: {"Content-Type": "application/json"},
    body
  })
  if (!res.ok) {
    alert(`creating channel err - ${await res.text()}`)
    return
  }
  return await res.json()
}

async function fetchChannels(userId) {
  const url = `http://localhost:3000/users/${userId}/channels`;
  const res = await fetch(url);
  if (!res.ok) {
    console.log("err -", await res.text());
    return;
  }
  const {channels} = await res.json();
  return channels
}

async function fetchChannelById(channelId) {
  const res = await fetch(`http://localhost:3000/channels/${channelId}`)
  if (!res.ok) {
    console.log('fetch channel by id err - ', await res.text())
    return
  }
  return await res.json()
}

export const channelsStore = createChannelsStore();
