# build stage
FROM node:lts-alpine AS build-stage

WORKDIR /app

# копируем package.json из папки frontend
COPY frontend/package*.json ./
RUN npm ci

# копируем все исходники фронтенда
COPY frontend/ ./

# собираем проект (должен быть скрипт "build" в frontend/package.json)
RUN npm run build

# production stage
FROM nginx:stable-alpine

# копируем собранные файлы
COPY --from=build-stage /app/dist /usr/share/nginx/html
# копируем конфиг nginx (если есть)
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]