Results of running
```
benchstat -html <(go test -count=5 -benchtime=0.2s -bench=././././serial) >results.md
```

The benchmark names follow the format `Adapter/(body_size)/(encoder)/(level)/serial-*`.

The first table lists time to compress the body with the specified encoder and level.

The second table lists the size of the resulting compressed body.

<table class='benchstat '>


<tbody>

<tr><th><th>time/op
<tr><td>Adapter/100/stdlib-gzip/1/serial-4<td>815ns ± 5%
<tr><td>Adapter/100/klauspost-gzip/1/serial-4<td>846ns ± 7%
<tr><td>Adapter/100/andybalholm-brotli/1/serial-4<td>788ns ± 2%
<tr><td>Adapter/100/google-cbrotli/1/serial-4<td>780ns ± 2%
<tr><td>Adapter/100/klauspost-zstd/1/serial-4<td>800ns ± 7%
<tr><td>Adapter/100/valyala-gozstd/1/serial-4<td>790ns ± 4%
<tr><td>Adapter/1000/google-cbrotli/1/serial-4<td>10.5µs ± 1%
<tr><td>Adapter/1000/google-cbrotli/2/serial-4<td>11.5µs ± 1%
<tr><td>Adapter/1000/google-cbrotli/3/serial-4<td>20.3µs ±16%
<tr><td>Adapter/1000/google-cbrotli/4/serial-4<td>29.9µs ± 2%
<tr><td>Adapter/1000/google-cbrotli/5/serial-4<td>34.4µs ± 1%
<tr><td>Adapter/1000/google-cbrotli/6/serial-4<td>35.3µs ± 3%
<tr><td>Adapter/1000/google-cbrotli/7/serial-4<td>44.2µs ±10%
<tr><td>Adapter/1000/google-cbrotli/8/serial-4<td>43.0µs ± 1%
<tr><td>Adapter/1000/google-cbrotli/9/serial-4<td>2.25ms ± 4%
<tr><td>Adapter/1000/google-cbrotli/10/serial-4<td>1.67ms ±11%
<tr><td>Adapter/1000/google-cbrotli/11/serial-4<td>2.58ms ± 3%
<tr><td>Adapter/1000/klauspost-zstd/1/serial-4<td>14.1µs ±30%
<tr><td>Adapter/1000/klauspost-zstd/2/serial-4<td>16.2µs ± 9%
<tr><td>Adapter/1000/klauspost-zstd/3/serial-4<td>25.8µs ± 9%
<tr><td>Adapter/1000/klauspost-zstd/4/serial-4<td>242µs ± 6%
<tr><td>Adapter/1000/valyala-gozstd/1/serial-4<td>10.3µs ± 5%
<tr><td>Adapter/1000/valyala-gozstd/2/serial-4<td>14.8µs ± 4%
<tr><td>Adapter/1000/valyala-gozstd/3/serial-4<td>32.2µs ±13%
<tr><td>Adapter/1000/valyala-gozstd/4/serial-4<td>76.5µs ± 2%
<tr><td>Adapter/1000/valyala-gozstd/5/serial-4<td>126µs ± 5%
<tr><td>Adapter/1000/valyala-gozstd/6/serial-4<td>124µs ± 1%
<tr><td>Adapter/1000/valyala-gozstd/7/serial-4<td>128µs ± 4%
<tr><td>Adapter/1000/valyala-gozstd/8/serial-4<td>128µs ± 5%
<tr><td>Adapter/1000/valyala-gozstd/9/serial-4<td>254µs ± 3%
<tr><td>Adapter/1000/valyala-gozstd/10/serial-4<td>513µs ± 1%
<tr><td>Adapter/1000/valyala-gozstd/11/serial-4<td>2.66ms ± 6%
<tr><td>Adapter/1000/valyala-gozstd/12/serial-4<td>2.64ms ± 9%
<tr><td>Adapter/1000/valyala-gozstd/13/serial-4<td>2.45ms ±11%
<tr><td>Adapter/1000/valyala-gozstd/14/serial-4<td>7.17ms ± 2%
<tr><td>Adapter/1000/valyala-gozstd/15/serial-4<td>9.35ms ± 1%
<tr><td>Adapter/1000/valyala-gozstd/16/serial-4<td>4.52ms ±13%
<tr><td>Adapter/1000/valyala-gozstd/17/serial-4<td>7.12ms ± 3%
<tr><td>Adapter/1000/valyala-gozstd/18/serial-4<td>7.20ms ± 2%
<tr><td>Adapter/1000/valyala-gozstd/19/serial-4<td>11.8ms ± 1%
<tr><td>Adapter/1000/valyala-gozstd/20/serial-4<td>23.6ms ± 2%
<tr><td>Adapter/1000/valyala-gozstd/21/serial-4<td>47.0ms ± 2%
<tr><td>Adapter/1000/valyala-gozstd/22/serial-4<td>95.5ms ± 1%
<tr><td>Adapter/1000/stdlib-gzip/1/serial-4<td>22.8µs ± 5%
<tr><td>Adapter/1000/stdlib-gzip/2/serial-4<td>58.7µs ± 3%
<tr><td>Adapter/1000/stdlib-gzip/3/serial-4<td>61.9µs ±18%
<tr><td>Adapter/1000/stdlib-gzip/4/serial-4<td>59.4µs ± 1%
<tr><td>Adapter/1000/stdlib-gzip/5/serial-4<td>62.4µs ± 4%
<tr><td>Adapter/1000/stdlib-gzip/6/serial-4<td>60.6µs ± 2%
<tr><td>Adapter/1000/stdlib-gzip/7/serial-4<td>61.8µs ± 2%
<tr><td>Adapter/1000/stdlib-gzip/8/serial-4<td>63.2µs ± 4%
<tr><td>Adapter/1000/stdlib-gzip/9/serial-4<td>63.3µs ± 7%
<tr><td>Adapter/1000/klauspost-gzip/1/serial-4<td>18.3µs ±11%
<tr><td>Adapter/1000/klauspost-gzip/2/serial-4<td>17.9µs ± 6%
<tr><td>Adapter/1000/klauspost-gzip/3/serial-4<td>18.8µs ± 8%
<tr><td>Adapter/1000/klauspost-gzip/4/serial-4<td>18.9µs ± 1%
<tr><td>Adapter/1000/klauspost-gzip/5/serial-4<td>21.2µs ± 5%
<tr><td>Adapter/1000/klauspost-gzip/6/serial-4<td>20.4µs ± 1%
<tr><td>Adapter/1000/klauspost-gzip/7/serial-4<td>43.8µs ± 3%
<tr><td>Adapter/1000/klauspost-gzip/8/serial-4<td>45.3µs ± 4%
<tr><td>Adapter/1000/klauspost-gzip/9/serial-4<td>46.4µs ± 7%
<tr><td>Adapter/1000/andybalholm-brotli/1/serial-4<td>23.0µs ± 1%
<tr><td>Adapter/1000/andybalholm-brotli/2/serial-4<td>30.1µs ± 3%
<tr><td>Adapter/1000/andybalholm-brotli/3/serial-4<td>51.0µs ± 2%
<tr><td>Adapter/1000/andybalholm-brotli/4/serial-4<td>71.5µs ± 3%
<tr><td>Adapter/1000/andybalholm-brotli/5/serial-4<td>73.7µs ± 3%
<tr><td>Adapter/1000/andybalholm-brotli/6/serial-4<td>74.8µs ±12%
<tr><td>Adapter/1000/andybalholm-brotli/7/serial-4<td>93.9µs ± 1%
<tr><td>Adapter/1000/andybalholm-brotli/8/serial-4<td>94.8µs ± 6%
<tr><td>Adapter/1000/andybalholm-brotli/9/serial-4<td>106µs ± 5%
<tr><td>Adapter/1000/andybalholm-brotli/10/serial-4<td>2.09ms ± 2%
<tr><td>Adapter/1000/andybalholm-brotli/11/serial-4<td>3.32ms ± 4%
<tr><td>Adapter/10000/klauspost-zstd/1/serial-4<td>57.7µs ±22%
<tr><td>Adapter/10000/klauspost-zstd/2/serial-4<td>85.0µs ± 5%
<tr><td>Adapter/10000/klauspost-zstd/3/serial-4<td>125µs ± 1%
<tr><td>Adapter/10000/klauspost-zstd/4/serial-4<td>1.00ms ± 7%
<tr><td>Adapter/10000/valyala-gozstd/1/serial-4<td>30.1µs ±11%
<tr><td>Adapter/10000/valyala-gozstd/2/serial-4<td>31.4µs ± 1%
<tr><td>Adapter/10000/valyala-gozstd/3/serial-4<td>52.0µs ± 5%
<tr><td>Adapter/10000/valyala-gozstd/4/serial-4<td>98.6µs ± 1%
<tr><td>Adapter/10000/valyala-gozstd/5/serial-4<td>179µs ± 7%
<tr><td>Adapter/10000/valyala-gozstd/6/serial-4<td>176µs ± 1%
<tr><td>Adapter/10000/valyala-gozstd/7/serial-4<td>205µs ± 4%
<tr><td>Adapter/10000/valyala-gozstd/8/serial-4<td>226µs ±15%
<tr><td>Adapter/10000/valyala-gozstd/9/serial-4<td>349µs ± 2%
<tr><td>Adapter/10000/valyala-gozstd/10/serial-4<td>616µs ± 4%
<tr><td>Adapter/10000/valyala-gozstd/11/serial-4<td>2.32ms ±15%
<tr><td>Adapter/10000/valyala-gozstd/12/serial-4<td>2.24ms ±13%
<tr><td>Adapter/10000/valyala-gozstd/13/serial-4<td>2.44ms ±15%
<tr><td>Adapter/10000/valyala-gozstd/14/serial-4<td>7.31ms ± 3%
<tr><td>Adapter/10000/valyala-gozstd/15/serial-4<td>9.61ms ± 1%
<tr><td>Adapter/10000/valyala-gozstd/16/serial-4<td>4.79ms ±11%
<tr><td>Adapter/10000/valyala-gozstd/17/serial-4<td>8.02ms ± 1%
<tr><td>Adapter/10000/valyala-gozstd/18/serial-4<td>8.79ms ± 3%
<tr><td>Adapter/10000/valyala-gozstd/19/serial-4<td>16.0ms ±15%
<tr><td>Adapter/10000/valyala-gozstd/20/serial-4<td>27.1ms ± 2%
<tr><td>Adapter/10000/valyala-gozstd/21/serial-4<td>50.6ms ± 3%
<tr><td>Adapter/10000/valyala-gozstd/22/serial-4<td>204ms ±87%
<tr><td>Adapter/10000/stdlib-gzip/1/serial-4<td>97.7µs ± 2%
<tr><td>Adapter/10000/stdlib-gzip/2/serial-4<td>148µs ±10%
<tr><td>Adapter/10000/stdlib-gzip/3/serial-4<td>144µs ± 1%
<tr><td>Adapter/10000/stdlib-gzip/4/serial-4<td>161µs ± 3%
<tr><td>Adapter/10000/stdlib-gzip/5/serial-4<td>177µs ± 4%
<tr><td>Adapter/10000/stdlib-gzip/6/serial-4<td>199µs ± 3%
<tr><td>Adapter/10000/stdlib-gzip/7/serial-4<td>251µs ± 2%
<tr><td>Adapter/10000/stdlib-gzip/8/serial-4<td>275µs ± 2%
<tr><td>Adapter/10000/stdlib-gzip/9/serial-4<td>284µs ± 2%
<tr><td>Adapter/10000/klauspost-gzip/1/serial-4<td>63.8µs ± 2%
<tr><td>Adapter/10000/klauspost-gzip/2/serial-4<td>64.4µs ± 5%
<tr><td>Adapter/10000/klauspost-gzip/3/serial-4<td>73.3µs ± 4%
<tr><td>Adapter/10000/klauspost-gzip/4/serial-4<td>86.0µs ± 6%
<tr><td>Adapter/10000/klauspost-gzip/5/serial-4<td>97.8µs ± 2%
<tr><td>Adapter/10000/klauspost-gzip/6/serial-4<td>102µs ± 4%
<tr><td>Adapter/10000/klauspost-gzip/7/serial-4<td>140µs ± 6%
<tr><td>Adapter/10000/klauspost-gzip/8/serial-4<td>165µs ± 2%
<tr><td>Adapter/10000/klauspost-gzip/9/serial-4<td>252µs ± 2%
<tr><td>Adapter/10000/andybalholm-brotli/1/serial-4<td>117µs ± 6%
<tr><td>Adapter/10000/andybalholm-brotli/2/serial-4<td>205µs ± 7%
<tr><td>Adapter/10000/andybalholm-brotli/3/serial-4<td>272µs ± 2%
<tr><td>Adapter/10000/andybalholm-brotli/4/serial-4<td>366µs ± 2%
<tr><td>Adapter/10000/andybalholm-brotli/5/serial-4<td>419µs ± 2%
<tr><td>Adapter/10000/andybalholm-brotli/6/serial-4<td>466µs ±13%
<tr><td>Adapter/10000/andybalholm-brotli/7/serial-4<td>579µs ± 3%
<tr><td>Adapter/10000/andybalholm-brotli/8/serial-4<td>616µs ± 2%
<tr><td>Adapter/10000/andybalholm-brotli/9/serial-4<td>718µs ± 1%
<tr><td>Adapter/10000/andybalholm-brotli/10/serial-4<td>8.20ms ± 5%
<tr><td>Adapter/10000/andybalholm-brotli/11/serial-4<td>21.8ms ± 3%
<tr><td>Adapter/10000/google-cbrotli/1/serial-4<td>36.6µs ±13%
<tr><td>Adapter/10000/google-cbrotli/2/serial-4<td>74.3µs ± 3%
<tr><td>Adapter/10000/google-cbrotli/3/serial-4<td>97.5µs ± 5%
<tr><td>Adapter/10000/google-cbrotli/4/serial-4<td>157µs ± 3%
<tr><td>Adapter/10000/google-cbrotli/5/serial-4<td>268µs ± 9%
<tr><td>Adapter/10000/google-cbrotli/6/serial-4<td>268µs ± 4%
<tr><td>Adapter/10000/google-cbrotli/7/serial-4<td>337µs ± 5%
<tr><td>Adapter/10000/google-cbrotli/8/serial-4<td>358µs ± 8%
<tr><td>Adapter/10000/google-cbrotli/9/serial-4<td>3.50ms ±140%
<tr><td>Adapter/10000/google-cbrotli/10/serial-4<td>5.60ms ± 3%
<tr><td>Adapter/10000/google-cbrotli/11/serial-4<td>15.1ms ± 2%
<tr><td>Adapter/100000/klauspost-zstd/1/serial-4<td>655µs ± 1%
<tr><td>Adapter/100000/klauspost-zstd/2/serial-4<td>819µs ± 8%
<tr><td>Adapter/100000/klauspost-zstd/3/serial-4<td>1.09ms ± 6%
<tr><td>Adapter/100000/klauspost-zstd/4/serial-4<td>6.69ms ±16%
<tr><td>Adapter/100000/valyala-gozstd/1/serial-4<td>371µs ± 2%
<tr><td>Adapter/100000/valyala-gozstd/2/serial-4<td>364µs ± 1%
<tr><td>Adapter/100000/valyala-gozstd/3/serial-4<td>512µs ± 1%
<tr><td>Adapter/100000/valyala-gozstd/4/serial-4<td>547µs ± 9%
<tr><td>Adapter/100000/valyala-gozstd/5/serial-4<td>994µs ±12%
<tr><td>Adapter/100000/valyala-gozstd/6/serial-4<td>982µs ± 6%
<tr><td>Adapter/100000/valyala-gozstd/7/serial-4<td>1.27ms ± 2%
<tr><td>Adapter/100000/valyala-gozstd/8/serial-4<td>1.44ms ± 5%
<tr><td>Adapter/100000/valyala-gozstd/9/serial-4<td>1.72ms ± 3%
<tr><td>Adapter/100000/valyala-gozstd/10/serial-4<td>2.26ms ±14%
<tr><td>Adapter/100000/valyala-gozstd/11/serial-4<td>3.93ms ± 7%
<tr><td>Adapter/100000/valyala-gozstd/12/serial-4<td>3.93ms ±12%
<tr><td>Adapter/100000/valyala-gozstd/13/serial-4<td>7.13ms ± 3%
<tr><td>Adapter/100000/valyala-gozstd/14/serial-4<td>11.9ms ± 2%
<tr><td>Adapter/100000/valyala-gozstd/15/serial-4<td>15.1ms ± 2%
<tr><td>Adapter/100000/valyala-gozstd/16/serial-4<td>19.4ms ± 2%
<tr><td>Adapter/100000/valyala-gozstd/17/serial-4<td>23.9ms ± 6%
<tr><td>Adapter/100000/valyala-gozstd/18/serial-4<td>31.0ms ± 1%
<tr><td>Adapter/100000/valyala-gozstd/19/serial-4<td>59.5ms ± 7%
<tr><td>Adapter/100000/valyala-gozstd/20/serial-4<td>74.0ms ± 9%
<tr><td>Adapter/100000/valyala-gozstd/21/serial-4<td>93.0ms ± 1%
<tr><td>Adapter/100000/valyala-gozstd/22/serial-4<td>145ms ± 6%
<tr><td>Adapter/100000/stdlib-gzip/1/serial-4<td>871µs ± 2%
<tr><td>Adapter/100000/stdlib-gzip/2/serial-4<td>1.04ms ± 4%
<tr><td>Adapter/100000/stdlib-gzip/3/serial-4<td>1.16ms ± 6%
<tr><td>Adapter/100000/stdlib-gzip/4/serial-4<td>1.30ms ± 2%
<tr><td>Adapter/100000/stdlib-gzip/5/serial-4<td>1.50ms ± 1%
<tr><td>Adapter/100000/stdlib-gzip/6/serial-4<td>1.71ms ± 3%
<tr><td>Adapter/100000/stdlib-gzip/7/serial-4<td>2.92ms ± 3%
<tr><td>Adapter/100000/stdlib-gzip/8/serial-4<td>7.20ms ± 7%
<tr><td>Adapter/100000/stdlib-gzip/9/serial-4<td>7.06ms ± 2%
<tr><td>Adapter/100000/klauspost-gzip/1/serial-4<td>617µs ± 3%
<tr><td>Adapter/100000/klauspost-gzip/2/serial-4<td>604µs ± 5%
<tr><td>Adapter/100000/klauspost-gzip/3/serial-4<td>792µs ± 5%
<tr><td>Adapter/100000/klauspost-gzip/4/serial-4<td>907µs ± 5%
<tr><td>Adapter/100000/klauspost-gzip/5/serial-4<td>1.06ms ± 2%
<tr><td>Adapter/100000/klauspost-gzip/6/serial-4<td>1.15ms ± 4%
<tr><td>Adapter/100000/klauspost-gzip/7/serial-4<td>1.30ms ± 3%
<tr><td>Adapter/100000/klauspost-gzip/8/serial-4<td>1.62ms ± 2%
<tr><td>Adapter/100000/klauspost-gzip/9/serial-4<td>6.78ms ± 5%
<tr><td>Adapter/100000/andybalholm-brotli/1/serial-4<td>1.14ms ± 3%
<tr><td>Adapter/100000/andybalholm-brotli/2/serial-4<td>2.06ms ±15%
<tr><td>Adapter/100000/andybalholm-brotli/3/serial-4<td>2.45ms ± 2%
<tr><td>Adapter/100000/andybalholm-brotli/4/serial-4<td>3.08ms ± 3%
<tr><td>Adapter/100000/andybalholm-brotli/5/serial-4<td>3.95ms ± 4%
<tr><td>Adapter/100000/andybalholm-brotli/6/serial-4<td>4.21ms ± 9%
<tr><td>Adapter/100000/andybalholm-brotli/7/serial-4<td>5.36ms ± 5%
<tr><td>Adapter/100000/andybalholm-brotli/8/serial-4<td>6.01ms ± 5%
<tr><td>Adapter/100000/andybalholm-brotli/9/serial-4<td>7.70ms ± 2%
<tr><td>Adapter/100000/andybalholm-brotli/10/serial-4<td>83.6ms ± 4%
<tr><td>Adapter/100000/andybalholm-brotli/11/serial-4<td>233ms ± 1%
<tr><td>Adapter/100000/google-cbrotli/1/serial-4<td>376µs ± 7%
<tr><td>Adapter/100000/google-cbrotli/2/serial-4<td>767µs ±11%
<tr><td>Adapter/100000/google-cbrotli/3/serial-4<td>834µs ± 0%
<tr><td>Adapter/100000/google-cbrotli/4/serial-4<td>1.17ms ± 3%
<tr><td>Adapter/100000/google-cbrotli/5/serial-4<td>2.39ms ± 3%
<tr><td>Adapter/100000/google-cbrotli/6/serial-4<td>2.49ms ± 3%
<tr><td>Adapter/100000/google-cbrotli/7/serial-4<td>3.09ms ± 2%
<tr><td>Adapter/100000/google-cbrotli/8/serial-4<td>3.41ms ± 2%
<tr><td>Adapter/100000/google-cbrotli/9/serial-4<td>8.66ms ±62%
<tr><td>Adapter/100000/google-cbrotli/10/serial-4<td>55.3ms ± 5%
<tr><td>Adapter/100000/google-cbrotli/11/serial-4<td>163ms ± 2%
<tr><td>&nbsp;
</tbody>

<tbody>

<tr><th><th>%
<tr><td>Adapter/100/stdlib-gzip/1/serial-4<td>100 ± 0%
<tr><td>Adapter/100/klauspost-gzip/1/serial-4<td>100 ± 0%
<tr><td>Adapter/100/andybalholm-brotli/1/serial-4<td>100 ± 0%
<tr><td>Adapter/100/google-cbrotli/1/serial-4<td>100 ± 0%
<tr><td>Adapter/100/klauspost-zstd/1/serial-4<td>100 ± 0%
<tr><td>Adapter/100/valyala-gozstd/1/serial-4<td>100 ± 0%
<tr><td>Adapter/1000/google-cbrotli/1/serial-4<td>45.2 ± 0%
<tr><td>Adapter/1000/google-cbrotli/2/serial-4<td>42.1 ± 0%
<tr><td>Adapter/1000/google-cbrotli/3/serial-4<td>39.9 ± 0%
<tr><td>Adapter/1000/google-cbrotli/4/serial-4<td>39.3 ± 0%
<tr><td>Adapter/1000/google-cbrotli/5/serial-4<td>36.6 ± 0%
<tr><td>Adapter/1000/google-cbrotli/6/serial-4<td>36.8 ± 0%
<tr><td>Adapter/1000/google-cbrotli/7/serial-4<td>36.7 ± 0%
<tr><td>Adapter/1000/google-cbrotli/8/serial-4<td>36.7 ± 0%
<tr><td>Adapter/1000/google-cbrotli/9/serial-4<td>36.7 ± 0%
<tr><td>Adapter/1000/google-cbrotli/10/serial-4<td>37.4 ± 0%
<tr><td>Adapter/1000/google-cbrotli/11/serial-4<td>37.2 ± 0%
<tr><td>Adapter/1000/klauspost-zstd/1/serial-4<td>42.9 ± 0%
<tr><td>Adapter/1000/klauspost-zstd/2/serial-4<td>42.2 ± 0%
<tr><td>Adapter/1000/klauspost-zstd/3/serial-4<td>41.3 ± 0%
<tr><td>Adapter/1000/klauspost-zstd/4/serial-4<td>40.5 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/1/serial-4<td>43.0 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/2/serial-4<td>42.6 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/3/serial-4<td>42.1 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/4/serial-4<td>42.1 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/5/serial-4<td>42.1 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/6/serial-4<td>41.3 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/7/serial-4<td>41.3 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/8/serial-4<td>41.3 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/9/serial-4<td>40.9 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/10/serial-4<td>40.9 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/11/serial-4<td>40.9 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/12/serial-4<td>40.9 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/13/serial-4<td>40.5 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/14/serial-4<td>40.5 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/15/serial-4<td>40.5 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/16/serial-4<td>40.4 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/17/serial-4<td>40.1 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/18/serial-4<td>40.0 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/19/serial-4<td>40.0 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/20/serial-4<td>40.0 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/21/serial-4<td>40.0 ± 0%
<tr><td>Adapter/1000/valyala-gozstd/22/serial-4<td>40.0 ± 0%
<tr><td>Adapter/1000/stdlib-gzip/1/serial-4<td>44.2 ± 0%
<tr><td>Adapter/1000/stdlib-gzip/2/serial-4<td>41.9 ± 0%
<tr><td>Adapter/1000/stdlib-gzip/3/serial-4<td>41.7 ± 0%
<tr><td>Adapter/1000/stdlib-gzip/4/serial-4<td>41.7 ± 0%
<tr><td>Adapter/1000/stdlib-gzip/5/serial-4<td>41.1 ± 0%
<tr><td>Adapter/1000/stdlib-gzip/6/serial-4<td>41.1 ± 0%
<tr><td>Adapter/1000/stdlib-gzip/7/serial-4<td>41.0 ± 0%
<tr><td>Adapter/1000/stdlib-gzip/8/serial-4<td>41.0 ± 0%
<tr><td>Adapter/1000/stdlib-gzip/9/serial-4<td>41.0 ± 0%
<tr><td>Adapter/1000/klauspost-gzip/1/serial-4<td>45.3 ± 0%
<tr><td>Adapter/1000/klauspost-gzip/2/serial-4<td>44.2 ± 0%
<tr><td>Adapter/1000/klauspost-gzip/3/serial-4<td>42.2 ± 0%
<tr><td>Adapter/1000/klauspost-gzip/4/serial-4<td>43.5 ± 0%
<tr><td>Adapter/1000/klauspost-gzip/5/serial-4<td>42.2 ± 0%
<tr><td>Adapter/1000/klauspost-gzip/6/serial-4<td>42.0 ± 0%
<tr><td>Adapter/1000/klauspost-gzip/7/serial-4<td>41.0 ± 0%
<tr><td>Adapter/1000/klauspost-gzip/8/serial-4<td>40.7 ± 0%
<tr><td>Adapter/1000/klauspost-gzip/9/serial-4<td>40.7 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/1/serial-4<td>45.2 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/2/serial-4<td>42.1 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/3/serial-4<td>39.9 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/4/serial-4<td>39.3 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/5/serial-4<td>36.6 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/6/serial-4<td>36.8 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/7/serial-4<td>36.7 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/8/serial-4<td>36.7 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/9/serial-4<td>36.7 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/10/serial-4<td>37.4 ± 0%
<tr><td>Adapter/1000/andybalholm-brotli/11/serial-4<td>37.2 ± 0%
<tr><td>Adapter/10000/klauspost-zstd/1/serial-4<td>28.6 ± 0%
<tr><td>Adapter/10000/klauspost-zstd/2/serial-4<td>28.1 ± 0%
<tr><td>Adapter/10000/klauspost-zstd/3/serial-4<td>27.6 ± 0%
<tr><td>Adapter/10000/klauspost-zstd/4/serial-4<td>27.5 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/1/serial-4<td>29.1 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/2/serial-4<td>28.5 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/3/serial-4<td>28.2 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/4/serial-4<td>28.2 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/5/serial-4<td>28.0 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/6/serial-4<td>27.6 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/7/serial-4<td>27.4 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/8/serial-4<td>27.3 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/9/serial-4<td>27.1 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/10/serial-4<td>27.1 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/11/serial-4<td>27.1 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/12/serial-4<td>27.0 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/13/serial-4<td>26.9 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/14/serial-4<td>26.9 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/15/serial-4<td>26.9 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/16/serial-4<td>26.8 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/17/serial-4<td>26.6 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/18/serial-4<td>26.5 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/19/serial-4<td>26.4 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/20/serial-4<td>26.4 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/21/serial-4<td>26.4 ± 0%
<tr><td>Adapter/10000/valyala-gozstd/22/serial-4<td>26.4 ± 0%
<tr><td>Adapter/10000/stdlib-gzip/1/serial-4<td>29.6 ± 0%
<tr><td>Adapter/10000/stdlib-gzip/2/serial-4<td>28.5 ± 0%
<tr><td>Adapter/10000/stdlib-gzip/3/serial-4<td>28.1 ± 0%
<tr><td>Adapter/10000/stdlib-gzip/4/serial-4<td>28.2 ± 0%
<tr><td>Adapter/10000/stdlib-gzip/5/serial-4<td>27.7 ± 0%
<tr><td>Adapter/10000/stdlib-gzip/6/serial-4<td>27.6 ± 0%
<tr><td>Adapter/10000/stdlib-gzip/7/serial-4<td>27.2 ± 0%
<tr><td>Adapter/10000/stdlib-gzip/8/serial-4<td>27.2 ± 0%
<tr><td>Adapter/10000/stdlib-gzip/9/serial-4<td>27.2 ± 0%
<tr><td>Adapter/10000/klauspost-gzip/1/serial-4<td>31.0 ± 0%
<tr><td>Adapter/10000/klauspost-gzip/2/serial-4<td>30.4 ± 0%
<tr><td>Adapter/10000/klauspost-gzip/3/serial-4<td>29.1 ± 0%
<tr><td>Adapter/10000/klauspost-gzip/4/serial-4<td>29.2 ± 0%
<tr><td>Adapter/10000/klauspost-gzip/5/serial-4<td>28.3 ± 0%
<tr><td>Adapter/10000/klauspost-gzip/6/serial-4<td>28.1 ± 0%
<tr><td>Adapter/10000/klauspost-gzip/7/serial-4<td>27.9 ± 0%
<tr><td>Adapter/10000/klauspost-gzip/8/serial-4<td>27.4 ± 0%
<tr><td>Adapter/10000/klauspost-gzip/9/serial-4<td>27.2 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/1/serial-4<td>29.9 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/2/serial-4<td>28.1 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/3/serial-4<td>27.9 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/4/serial-4<td>27.5 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/5/serial-4<td>26.0 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/6/serial-4<td>25.9 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/7/serial-4<td>25.8 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/8/serial-4<td>25.8 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/9/serial-4<td>25.8 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/10/serial-4<td>23.4 ± 0%
<tr><td>Adapter/10000/andybalholm-brotli/11/serial-4<td>23.1 ± 0%
<tr><td>Adapter/10000/google-cbrotli/1/serial-4<td>29.9 ± 0%
<tr><td>Adapter/10000/google-cbrotli/2/serial-4<td>28.1 ± 0%
<tr><td>Adapter/10000/google-cbrotli/3/serial-4<td>27.9 ± 0%
<tr><td>Adapter/10000/google-cbrotli/4/serial-4<td>27.5 ± 0%
<tr><td>Adapter/10000/google-cbrotli/5/serial-4<td>26.0 ± 0%
<tr><td>Adapter/10000/google-cbrotli/6/serial-4<td>25.9 ± 0%
<tr><td>Adapter/10000/google-cbrotli/7/serial-4<td>25.8 ± 0%
<tr><td>Adapter/10000/google-cbrotli/8/serial-4<td>25.8 ± 0%
<tr><td>Adapter/10000/google-cbrotli/9/serial-4<td>25.8 ± 0%
<tr><td>Adapter/10000/google-cbrotli/10/serial-4<td>23.4 ± 0%
<tr><td>Adapter/10000/google-cbrotli/11/serial-4<td>23.1 ± 0%
<tr><td>Adapter/100000/klauspost-zstd/1/serial-4<td>25.9 ± 0%
<tr><td>Adapter/100000/klauspost-zstd/2/serial-4<td>25.8 ± 0%
<tr><td>Adapter/100000/klauspost-zstd/3/serial-4<td>25.5 ± 0%
<tr><td>Adapter/100000/klauspost-zstd/4/serial-4<td>24.7 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/1/serial-4<td>26.4 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/2/serial-4<td>25.8 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/3/serial-4<td>25.8 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/4/serial-4<td>25.8 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/5/serial-4<td>25.2 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/6/serial-4<td>25.0 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/7/serial-4<td>24.5 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/8/serial-4<td>24.2 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/9/serial-4<td>24.1 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/10/serial-4<td>24.1 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/11/serial-4<td>24.1 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/12/serial-4<td>24.0 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/13/serial-4<td>23.6 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/14/serial-4<td>23.6 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/15/serial-4<td>23.6 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/16/serial-4<td>23.2 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/17/serial-4<td>23.0 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/18/serial-4<td>22.9 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/19/serial-4<td>22.7 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/20/serial-4<td>22.7 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/21/serial-4<td>22.7 ± 0%
<tr><td>Adapter/100000/valyala-gozstd/22/serial-4<td>22.7 ± 0%
<tr><td>Adapter/100000/stdlib-gzip/1/serial-4<td>27.3 ± 0%
<tr><td>Adapter/100000/stdlib-gzip/2/serial-4<td>26.2 ± 0%
<tr><td>Adapter/100000/stdlib-gzip/3/serial-4<td>25.6 ± 0%
<tr><td>Adapter/100000/stdlib-gzip/4/serial-4<td>25.8 ± 0%
<tr><td>Adapter/100000/stdlib-gzip/5/serial-4<td>24.9 ± 0%
<tr><td>Adapter/100000/stdlib-gzip/6/serial-4<td>24.8 ± 0%
<tr><td>Adapter/100000/stdlib-gzip/7/serial-4<td>24.2 ± 0%
<tr><td>Adapter/100000/stdlib-gzip/8/serial-4<td>24.1 ± 0%
<tr><td>Adapter/100000/stdlib-gzip/9/serial-4<td>24.1 ± 0%
<tr><td>Adapter/100000/klauspost-gzip/1/serial-4<td>28.5 ± 0%
<tr><td>Adapter/100000/klauspost-gzip/2/serial-4<td>27.6 ± 0%
<tr><td>Adapter/100000/klauspost-gzip/3/serial-4<td>26.9 ± 0%
<tr><td>Adapter/100000/klauspost-gzip/4/serial-4<td>26.4 ± 0%
<tr><td>Adapter/100000/klauspost-gzip/5/serial-4<td>25.8 ± 0%
<tr><td>Adapter/100000/klauspost-gzip/6/serial-4<td>25.7 ± 0%
<tr><td>Adapter/100000/klauspost-gzip/7/serial-4<td>25.4 ± 0%
<tr><td>Adapter/100000/klauspost-gzip/8/serial-4<td>24.6 ± 0%
<tr><td>Adapter/100000/klauspost-gzip/9/serial-4<td>24.1 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/1/serial-4<td>27.6 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/2/serial-4<td>25.4 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/3/serial-4<td>25.3 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/4/serial-4<td>24.9 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/5/serial-4<td>23.6 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/6/serial-4<td>23.4 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/7/serial-4<td>23.3 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/8/serial-4<td>23.2 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/9/serial-4<td>23.2 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/10/serial-4<td>20.2 ± 0%
<tr><td>Adapter/100000/andybalholm-brotli/11/serial-4<td>19.9 ± 0%
<tr><td>Adapter/100000/google-cbrotli/1/serial-4<td>27.6 ± 0%
<tr><td>Adapter/100000/google-cbrotli/2/serial-4<td>25.4 ± 0%
<tr><td>Adapter/100000/google-cbrotli/3/serial-4<td>25.3 ± 0%
<tr><td>Adapter/100000/google-cbrotli/4/serial-4<td>24.9 ± 0%
<tr><td>Adapter/100000/google-cbrotli/5/serial-4<td>23.6 ± 0%
<tr><td>Adapter/100000/google-cbrotli/6/serial-4<td>23.4 ± 0%
<tr><td>Adapter/100000/google-cbrotli/7/serial-4<td>23.3 ± 0%
<tr><td>Adapter/100000/google-cbrotli/8/serial-4<td>23.2 ± 0%
<tr><td>Adapter/100000/google-cbrotli/9/serial-4<td>23.2 ± 0%
<tr><td>Adapter/100000/google-cbrotli/10/serial-4<td>20.2 ± 0%
<tr><td>Adapter/100000/google-cbrotli/11/serial-4<td>19.9 ± 0%
<tr><td>&nbsp;
</tbody>

</table>
</body>
</html>
