#!/bin/bash

for ((i=1;i<=$(ls -1q testdata/test*.txt | wc -l);i++)) ; do
  f="testdata/test$i.txt"
  output="testdata/out$i.txt"
  answers="testdata/answers$i.txt"
  printf "Processing test $i \n";
  ./medasync-challenge -f $f > $output
  if diff $output $answers; then
    printf "Passed test $i\n"
  else
    printf "Failed test $i\n"
  fi
done
