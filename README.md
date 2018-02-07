# hashCode2017
My attempt at google's hashCode 2017 Final

[The Problem](https://hashcode.withgoogle.com/2017/tasks/hashcode2017_final_task.pdf)

## Running it
**Linux**

```./hashCode2017 <input-filename> <output-filename>```

**Windows**

```hashCode2017.exe <input-filename> <output-filename>```


Files can be found in *final_round_2017.in*

## How it works
1. Parse file
2. Loop though grid until fully covered or budget exceeded
    1. If max rating found add router
    2. Add router to backbone
3. Print everything

## Make it better
* Backbone mapping can be diagonal
* Make it faster somehow
* Make the rating function better to adjust for distance from backbone
