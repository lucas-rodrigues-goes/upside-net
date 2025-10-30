minifab up -i 1.4.4 -o hawkins.com
minifab ccup -v 3.0 -n blendchaincode  -l go
minifab instantiate -n blendchaincode -p '' -c mychannel
minifab invoke -p \"initDimensionalEnergy\",\"1\",\"lab\",\"55\",\"2\",\"Lucas\" -n blendchaincode