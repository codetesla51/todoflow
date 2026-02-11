FROM node:20-slim

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

# Environment variable for the API URL
ENV PUBLIC_API_URL=http://localhost:8080

EXPOSE 5173

CMD ["npm", "run", "dev", "--", "--host"]
