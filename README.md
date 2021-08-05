# ReApp module for assigning students to workplaces
This module is responsible for identification of optimal distribution students to workplaces where they will do apprentice. It is GO code what can be compiled localy, but primary should be used in GCP Function FaaS environment. 

## Used methods
Optimization algorithm is Artifical Bee Colony without any additional modification. There are two criterias used for distribution (fitness calculation) distance between workplace and home (used only city names not a specific addresses) and if specific student was already on specific workplace before.

## Development notices
All development actions are coordinated via Makefile

### Prerequisites
1. You have to have account by Google GCP with enabled GCP Functions and Google Maps DistanceMatrix API
1. DistanceMatrixAPI key store localy in file `secret`. Compile through `make local` or `make install` will run sed to replace PLACEHOLDER in distance.go and then replace back PLACEHOLDER (not to leak token to GitHub) 

### Local development
Please use `make local` for tests execution and compile code. Resulting binary is stored in bin/ directory.

### FaaS deployment
Use `make install`, what executes tests and gcp cli tool to upload to GCP
