# Use uma imagem base com suporte para Go
FROM golang:1.22

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copie os arquivos do projeto para o contêiner
COPY . .

# Baixe as dependências
RUN go mod tidy

# Compile o binário
RUN go build -o main .

# Exponha a porta da API
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./main"]
