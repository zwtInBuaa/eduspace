# 基础镜像
FROM node:lts-alpine AS build
# 工作目录
RUN mkdir /build
WORKDIR /build
RUN npm config set registry https://registry.npm.taobao.org
# 安装依赖
COPY ["./package.json", "./package-lock.json*", "./npm-shrinkwrap.json*", "./"]
RUN npm install
# 将上层目录中的所有文件复制到容器的 WORKDIR 目录中
COPY . .
# 构建应用程序
RUN npm run build

FROM node:lts-alpine AS runtime
# 工作目录
RUN mkdir /app
WORKDIR /app
RUN npm install -g serve
COPY --from=build /build/dist /app/dist
# 构建应用程序
EXPOSE 6060
CMD ["serve", "dist", "-p", "6060"]

