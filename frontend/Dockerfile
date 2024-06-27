FROM node:latest

RUN npm install -g @angular/cli

WORKDIR /app

COPY . .

RUN npm install
CMD ["ng", "serve", "--host", "0.0.0.0"]



