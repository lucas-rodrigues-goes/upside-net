# UpsideNet Chaincode - Proof of Concept

This chaincode manages two types of data:
1. **Dimensional Energy Measurements** - Stored in a PRIVATE collection (only accessible by Hawkins and Montauk)
2. **Environmental Control Data** - Stored in PUBLIC state (accessible by all organizations)

## Setup

The chaincode has been created in `vars/chaincode/blendchaincode/go/`

## Deployment with Minifabric

To install and instantiate this chaincode on your minifabric network:

```bash
# Install the chaincode on all peers
minifab ccup -n blendchaincode -p vars/chaincode/blendchaincode/go -l go -v 1.0

# Or manually:
# minifab install -n blendchaincode
# minifab instantiate -n blendchaincode -c vars/chaincode/blendchaincode/collection_config.json
```

## Testing the Chaincode

### 1. Create a Dimensional Energy Measurement (PRIVATE - only Hawkins and Montauk can see)

```bash
minifab invoke -n blendchaincode -p '"initDimensionalEnergy","dim1","labA","1250.5","440.2","HawkinsLab"'
```

### 2. Create Environmental Control Data (PUBLIC - all organizations can see)

```bash
minifab invoke -n blendchaincode -p '"initEnvironmentalControl","env1","labA","22.5","65.0","1013.25","MontaukLab"'
```

### 3. Query Dimensional Energy Measurement

```bash
minifab invoke -n blendchaincode -p '"readDimensionalEnergy","dim1"'
```

### 4. Query Environmental Control Data

```bash
minifab invoke -n blendchaincode -p '"readEnvironmentalControl","env1"'
```

### 5. Get All Dimensional Energy Measurements

```bash
minifab invoke -n blendchaincode -p '"getAllDimensionalEnergy"'
```

### 6. Get All Environmental Control Data

```bash
minifab invoke -n blendchaincode -p '"getAllEnvironmentalControl"'
```

## Verification

- Oak Ridge will NOT be able to read dimensional energy measurements (they're in a private collection)
- Oak Ridge CAN read environmental control data (stored in public state)
- All organizations (Hawkins, Montauk, Oak Ridge) can read environmental control data

