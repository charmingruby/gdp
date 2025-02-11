1) Entrega ordenada para aplicação baseado na ordem dos pacotes (ter # de sequência).

2) Confirmação acumulativa (ACK acumulativo) do destinatário para o remetente.

3) Adicione no protocolo um controle de fluxo, onde o remetente deve saber qual o tamanho da janela do destinatário, a fim de não afogá-lo.

4) Agora crie uma equação de controle de congestionamento, a fim de que, se a rede estiver apresentando perda (muitos pacotes com ACK pendentes, ACK duplicados ou timeout), ele deve ser utilizado para reduzir o fluxo de envio de pacotes. 
    - 4.1 Você deve propor um controle de congestionamento, que pode ser baseado em algum existente no TCP, no QUIC ou qualquer outro protocolo. 
    - 4.2 Lembre da Aula 13, onde há controle de congestionamento no TCP que utiliza uma janela "cwnd" e um variável "ssthresh" para controle das fases de "Slow Start" e "Congestiona Avoidance".

5) Avalie seu protocolo sobre 1 remetente (cliente) que envia um arquivo (ou dados sintéticos que preencham o payload) para 1 destinatário (servidor). 
    - 5.1 Esses dados devem ser equivalente a, pelo menos, 10.000 pacotes. 
    - 5.2 Para avaliar o controle de congestionamento, insira perdas arbitrárias (ou utilizando uma função rand()) de pacotes no destinatário (você pode fazer isso sorteando a cada chegada de um novo pacote se ele será contabilizado e processado ou descartado).