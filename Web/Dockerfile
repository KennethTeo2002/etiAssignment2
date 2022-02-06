# syntax=docker/dockerfile:1
FROM node:16

# Set destination for COPY
WORKDIR /frontend

COPY package*.json ./

RUN npm install

COPY . .
COPY ./public ./public

EXPOSE 8070

CMD [ "node", "main.js" ]