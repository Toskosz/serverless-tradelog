
# serverless tradelog

A serverless REST API built using the AWS-SDK for golang. It has some jwt authentication and CRUD functionalities for the tradelog's. I made this with the intention to have my trading diary available anywhere at any time.

## Diagram

![Diagram](https://github.com/Toskosz/serverless-tradelog/blob/main/media/serverless.jpg)

## tradelog

Tradelogs are the logging of the basic data in a day-trade operation in the B3 futures market.

Data | Type | Meaning 
--- | --- | --- 
Abertura | string | Datetime of the start of the operation 
Fechamento | string | Datetime of the end of the operation  
Ativo | string | Asset of the operation 
Resultado | float32 | Result of the operation 
Contratos | int | Number of contracts aka leverage 
MEP | float32 | Maximun positive exposure during operation 
MEN | float32 | Maximun negative exposure during operation 
TempoOperacaoSegundos | int | Duration of op. in seconds 
PrecoCompra | float32 | Buying price 
PrecoVenda | float32 | Selling price 
Revisado | bool | If the op. has been reviewed 
Desc | string | Description of the operation


## Environment Variables

To run this project, you will need to add the following environment variables to your lambda environment

`AWS_REGION`

`JWT_SECRET` 

For best practices purposes `JWT_SECRET`  should be in stored in the secrets manager, but in the spirit of cutting costs it's stored as a environment variable.


## API Reference

All calls need to have a API key in the header. 

#### Register new user

```http
  POST /api/register
```

| Parameter | Type 
| :-------- | :------- 
| `email` | `string` 
| `password` | `string` 


#### login
```http
  POST /api/login
```

| Parameter | Type     
| :-------- | :------- 
| `email` | `string` 
| `password` | `string` 

#### logout

```http
  POST /api/logout
```

#### Get specific log
```http
  GET /api/my/logs/?log-abertura=YYYY-MM-DDTHH:MM:SSZ
```

| Parameter | Type | Description
| :-------- | :------- | :-------
| `log-abertura` | `string` | Datetime of the start of the operation

#### Get all user logs
```http
  GET /api/my/logs
```

#### Create log
```http
  POST /api/my/logs
```

The following parameters must be in the request body:

| Parameter | Type 
| :-------- | :------- 
| `abertura` | `string` 
| `fechamento` | `string` 
| `ativo` | `string` 
| `resultado` | `float` 
| `contratos` | `int` 
| `mep` | `float` 
| `men` | `float` 
| `duracao`| `int` 
| `preco-compra` | `float` 
| `preco-venda`| `float` 
| `revisado` | `bool` 
| `descricao` | `string` 

#### Update log
```http
  PUT /api/my/logs
```

The following parameters must be in the request body:

| Parameter | Type 
| :-------- | :------- 
| `abertura` | `string` 
| `revisado` | `bool` 
| `descricao` | `string` 

#### Delete specific log
```http
  DELETE /api/my/logs/?log-abertura=YYYY-MM-DDTHH:MM:SSZ
```

| Parameter | Type | Description
| :-------- | :------- | :-------
| `log-abertura` | `string` | Datetime of the start of the operation
