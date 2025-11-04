# Notas

Caso o código vá a ser executado em uma máquina **Windows**, comandos que incluem `"` (aspas) devem incluir uma `\` (contra-barra) precedendo cada uma das aspas e a utilização do CMD para rodar os comandos. Em máquinas Linux não há essas necessidades.

---

## Inicialização

* **Incluir dependencias GO localmente**

  ```bash
  cd ./chaincode/
  go mod vendor
  cd ..
  ```

* **Inicializar Fabric (sem pasta vars ou renomeada)**

  ```bash
  minifab up -i 1.4.4 -c channel -o hawkins.com
  ```

* **Inicializar Fabric (pasta vars já criada)**

  ```bash
  minifab up -i 1.4.4 -c channel -l go -v 0.0 -r true -n upsidenet -o hawkins.com
  ```

* **Inicializar Chaincode (requer pasta vars com chaincode e collection)**

  ```bash
  minifab ccup -n upsidenet -l go -v 1.0 -r true -o hawkins.com
  ```

* **Realizar Discover do Chaincode como montauk.com**

  ```bash
  minifab discover -o montauk.com
  ```

---

## Testes de manipulação de energia (private collection)

1. **Criar novo registro de energia como Hawkins**

   ```bash
   minifab invoke -p "initDimensionalEnergy","lab","55","2","Lucas" -n upsidenet -o hawkins.com
   ```

2. **Tentar ler dados com todas as organizações**

   ```bash
   minifab invoke -p "getAllDimensionalEnergy" -n upsidenet -o hawkins.com
   minifab invoke -p "getAllDimensionalEnergy" -n upsidenet -o montauk.com
   minifab invoke -p "getAllDimensionalEnergy" -n upsidenet -o oakridge.com
   ```

3. **Tentar inserir novo registro de energia como Oakridge**

   ```bash
   minifab invoke -p "initDimensionalEnergy","lab","98","2","Pedro" -n upsidenet -o oakridge.com
   ```

4. **Verificar registro de energia inserido como Hawkins**

   ```bash
   minifab invoke -p "getAllDimensionalEnergy" -n upsidenet -o hawkins.com
   ```

---

## Testes de registros ambientais

1. **Criar registro de controle ambiental como Oakridge**

   ```bash
   minifab invoke -p "initEnvironmentalControl","outsidePark","28","49","10","Pedro" -n upsidenet -o oakridge.com
   ```

2. **Tentar ler dados com todas as organizações**

   ```bash
   minifab invoke -p "getAllEnvironmentalControl" -n upsidenet -o oakridge.com
   minifab invoke -p "getAllEnvironmentalControl" -n upsidenet -o montauk.com
   minifab invoke -p "getAllEnvironmentalControl" -n upsidenet -o hawkins.com
   ```