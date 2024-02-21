# Prepare node/npm
FROM node:20.11.1 AS node_base

COPY . ./app

RUN npm install

RUN npm run build

EXPOSE 8000

CMD ["npm", "run", "dev"]