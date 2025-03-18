# 基础镜像
FROM node:lts-alpine AS build
# 工作目录
RUN mkdir /app
WORKDIR /app
# 安装依赖
COPY ["package.json", "package-lock.json*", "npm-shrinkwrap.json*", "./"]
RUN npm install
# 将当前目录中的所有文件复制到容器的 WORKDIR 目录中
COPY . .
# 构建应用程序
EXPOSE 8686
# 启动应用程序
RUN npm run lint -- --no-fix
CMD ["npm", "run", "serve"]

