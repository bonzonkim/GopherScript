#!/bin/bash

add_numbers(){
  echo $(( $1 + $2 ))
}

result=$(add_numbers 2 3)
echo $result
