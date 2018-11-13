# genetic-algorithms-in-parallel-programming

## GA implementation
This is the implementation of a word-finding GA. The purpose of a GA always is to find a solution for a certain problem from a huge area of possible solutions. The implementation of this GA just ties to find a certain string for example your name, your favorites car’s name or a Shakespeare poem. Provided a Alphabet "a-z" or "a-zA-Z " it uses as search area and knowing the length of  the string you are searching for.

The implementation of the GA itself is quiet simple but this project focuses on parallelization and measuring the effects of certain parallelization possibilitys.

## Parallelization
It can take much time for GA to find a solution for a certain problem. But certain parts of a GA are simple to parallelize and it's also simple to run multiple GA at once and let them share their population to speed them up. Using Golang and it's build in features helps to parallelize  those parts. This implementation tries to focus on both parallelizing the GA and generating a GA network.

## Run GA in a network
There are multiple ways to build a GA network. Some of the more common are a ring network, a classic network where each node has 4 neighbors or even a hypercubes. This implementation gives you the possibility to setup the connection between GAs the way you want. You can build networks providing just an array. Each index in the first dimension and it's partner in the second dimension build a connection. Note that each connection goes ONE way and a GA could even connect to itself. The following array will setup a ring:
```golang
//four GAs in a ring each connected to it's neighbors 
ring4GA := [][]int{
	{0, 1, 0, 1},
	{1, 0, 1, 0},
	{0, 1, 0, 1},
	{1, 0, 1, 0},
}
```

Also you can provide the population size, mutation rate, alphabet to be searched which are the same for all GA in the network 

## Parallelize the GA itself 
After the network parallelization i will continue to parallelize the GA's fitness, crossover and mutation function.

## Measuring
To measure the quality of a certain GA setup we have the following values for each GA: 

 - Duration to find a solution
 - Average fitness of the population
 - Standard derivation from average fitness

To provide a deterministic behavior a GA and a GA network can always be started with a certain seed (otherwise it will use a random one). This will always lead to the same start population and often to the very same behavior in the whole evolutionary process. This way we can measure the GA's quality with different setups quiet reliable. 
