# Full Cycle OTEL
Projeto para atender o desafio de labs da full cycle sobre OTEL  
Objetivo: Desenvolver um sistema em Go que receba um CEP, identifica a cidade e retorna o clima atual (temperatura em graus celsius, fahrenheit e kelvin) juntamente com a cidade. Esse sistema deverá implementar OTEL(Open Telemetry) e Zipkin.

Basedo no cenário conhecido "Sistema de temperatura por CEP" denominado Serviço B, será incluso um novo projeto, denominado Serviço A.

Requisitos - Serviço A (responsável pelo input):

O sistema deve receber um input de 8 dígitos via POST, através do schema:  { "cep": "29902555" }  
O sistema deve validar se o input é valido (contem 8 dígitos) e é uma STRING  
Caso seja válido, será encaminhado para o Serviço B via HTTP  
Caso não seja válido, deve retornar:  
Código HTTP: 422  
Mensagem: invalid zipcode  
Requisitos - Serviço B (responsável pela orquestração):  

O sistema deve receber um CEP válido de 8 digitos  
O sistema deve realizar a pesquisa do CEP e encontrar o nome da localização, a partir disso, deverá retornar as temperaturas e formata-lás em: Celsius, Fahrenheit, Kelvin juntamente com o nome da localização.  
O sistema deve responder adequadamente nos seguintes cenários:  
Em caso de sucesso:  
Código HTTP: 200
Response Body: { "city: "São Paulo", "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }  
Em caso de falha, caso o CEP não seja válido (com formato correto):  
Código HTTP: 422  
Mensagem: invalid zipcode  
Em caso de falha, caso o CEP não seja encontrado:  
Código HTTP: 404  
Mensagem: can not find zipcode  
Após a implementação dos serviços, adicione a implementação do OTEL + Zipkin:

Implementar tracing distribuído entre Serviço A - Serviço B  
Utilizar span para medir o tempo de resposta do serviço de busca de CEP e busca de temperatura  

Dicas:
Utilize a API viaCEP (ou similar) para encontrar a localização que deseja consultar a temperatura: https://viacep.com.br/  
Utilize a API WeatherAPI (ou similar) para consultar as temperaturas desejadas: https://www.weatherapi.com/  
Para realizar a conversão de Celsius para Fahrenheit, utilize a seguinte fórmula: F = C * 1,8 + 32  
Para realizar a conversão de Celsius para Kelvin, utilize a seguinte fórmula: K = C + 273  
Sendo F = Fahrenheit  
Sendo C = Celsius  
Sendo K = Kelvin  
Para dúvidas da implementação do OTEL, você pode clicar aqui  
Para implementação de spans, você pode clicar aqui  
Você precisará utilizar um serviço de collector do OTEL  
Para mais informações sobre Zipkin, você pode clicar aqui  

Entrega:
O código-fonte completo da implementação.  
Documentação explicando como rodar o projeto em ambiente dev.  
Utilize docker/docker-compose para que possamos realizar os testes de sua aplicação.  

## Como rodar a aplicação

### 1 - Alterando a chave da api de clima  
  
Alterar no arquivo .env do projeto service-b-climate a váriavel KEY_CLIMATE trocando o valor {{YOUR KEY HERE}} pela 
chave gerada no https://www.weatherapi.com/

### 2 - Criando e executando as imagens no docker  
  
É preciso rodar o comando  
<b>docker-compose up -d</b>  
para criar as imagens e containers no docker


#### 3 - URL da API
  
Abaixo o comando curl de consumo da api  
<b>curl --location 'localhost:8081/cep/82720540'</b>  
substituir o cep no valor 82720540 pelo desejado para realizar a busca e visualizar as informações de clima

#### 4 - Visualizar as informações de tempo no zipkin
  
Abrir a URL http://localhost:9411/zipkin/ e clicar no botão RUN QUERY para observar as informações de tempo 
no servico-a-cep

