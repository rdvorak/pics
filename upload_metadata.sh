for YEAR
do
exiftool -J -my $( find $HOME/Pictures/$YEAR -type d -name "web") | curl -POST --data-binary @- http://localhost:8081/submit/metadata
done
