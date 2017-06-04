# csvmask
Mask or remove fields from a CSV file

## Usage

```
Usage of ./csvmask:
  -d value
        comma separated list of columns to remove (remove takes precidence over mask)
  -header
        indicate whether the csv containers a header row
  -i string
        input file (csv formatted)
  -m value
        comma separated list of columns to mask
  -o string
        output file (csv formatted)
```

## Example

```
./csvmask -i in.csv -d 0,1 -m 2 -header
first_name,last_name,email,comment
,,04a05b938ae4aaa403702e36ca3e47c65f22fc151d5179f596659a7d85b4e0ef,我说中文
,,2bc6d3e1b63d25102fec3a2a83fb4d614d61d2a879ad592c35f792e4164d2eb4,I speak English
,,a9a9f311a8374479539f982c6e7d0e3e9afc20c245a8ef83b3726434716a57c3,I don't speak
```
