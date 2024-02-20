# Prepare node/npm
FROM node:latest AS node_base

WORKDIR /app

COPY ./ ./

RUN npm install

RUN npm run build

EXPOSE 8080

CMD ["npm", "run", "build"]