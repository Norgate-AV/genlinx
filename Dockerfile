FROM node:lts-alpine3.15

RUN apk update && apk upgrade
RUN apk add git

RUN yarn global add genlinx
RUN chown -R node:node /usr/local/lib/node_modules

USER node
ENTRYPOINT [ "genlinx" ]
