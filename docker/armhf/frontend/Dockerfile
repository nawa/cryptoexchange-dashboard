FROM node:9 as builder

RUN mkdir -p /usr/src/app/build
WORKDIR /usr/src/app
ENV PATH /usr/src/app/node_modules/.bin:$PATH

COPY package.json /usr/src/app/package.json

RUN npm install --silent
RUN npm install react-scripts@1.1.4 -g --silent

COPY . /usr/src/app

RUN REACT_APP_STAGE=prod npm run build

FROM arm32v7/nginx:1.14
COPY --from=builder /usr/src/app/build /usr/share/nginx/html
EXPOSE 80
VOLUME /usr/src/app/
CMD ["nginx", "-g", "daemon off;"]