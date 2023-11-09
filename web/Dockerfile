FROM node:20.5.1 as deps

WORKDIR /app
COPY package.json /app/package.json
COPY yarn.lock /app/yarn.lock
RUN yarn

FROM deps as builder
COPY . .
RUN yarn build

FROM node:20-alpine
WORKDIR /app

COPY --from=builder /app/dist dist
COPY --from=builder /app/cleardata.html cleardata.html
COPY --from=builder /app/server.js server.js
COPY package.json package.json
RUN npm install --omit=dev

CMD PORT=8090 node server.js

