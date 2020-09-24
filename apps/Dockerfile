############# BUILD STAGE ##############
FROM node:14.8.0-alpine as build
ARG APP_SRC

RUN apk add --no-cache --virtual .gyp \
        python \
        make \
        g++
RUN npm install -g babel-cli

USER node
WORKDIR /home/node

# Client app: stage package files
RUN mkdir ./client
COPY ./apps-container/client/package.json ./client/package.json
COPY ./apps-container/client/package-lock.json ./client/package-lock.json

# Server app: stage package files
RUN mkdir ./server
COPY ./apps-container/server/package.json ./server/package.json
COPY ./apps-container/server/package-lock.json ./server/package-lock.json

# Install deps
RUN cd client && npm install
RUN cd server && npm install

# Build for production
COPY --chown=node:node ./apps-container .

# Copy application files
COPY $APP_SRC/ui ./client/src/HurraApp
COPY $APP_SRC/backend ./server/src/HurraApp

RUN cd ./client && npm run-script build && npm upgrade caniuse-lite browserslist
RUN cd ./server && npm run-script build && npm upgrade caniuse-lite browserslist



######### RUN-TIME IMAGE ##############
FROM node:14.8.0-alpine
RUN npm install -g serve concurrently

USER node
WORKDIR /home/node
COPY --from=build /home/node/client/build ./client
COPY --from=build /home/node/server/dist ./server

CMD ["concurrently", "--kill-others-on-fail", "cd server && node server.js HurraServer.js", "cd client && HURRA_HOSTNAME=$HURRA_HOSTNAME serve -p 3000"]
