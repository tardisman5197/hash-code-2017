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
* Map backbone after router placement
* Make it faster somehow
* Make the rating function better to adjust for distance from backbone

## Scores
2017 HashCode [Scoreboard](https://hashcode.withgoogle.com/hashcode_2017.html)

My Overall Score: 539,378,943
### Files
[test.in](https://github.com/tardisman5197/hashCode2017/blob/master/final_round_2017.in/test.in): 65,907

[charleston_road.in](https://github.com/tardisman5197/hashCode2017/blob/master/final_round_2017.in/charleston_road.in): 21,959,496

[rue_de_londres.in](https://github.com/tardisman5197/hashCode2017/blob/master/final_round_2017.in/rue_de_londres.in): 57,509,784

[opera.in](https://github.com/tardisman5197/hashCode2017/blob/master/final_round_2017.in/opera.in): 169,719,814

[lets_go_higher.in](https://github.com/tardisman5197/hashCode2017/blob/master/final_round_2017.in/lets_go_higher.in): 290,123,942
