
## ref


https://docs.aws.amazon.com/aws-cost-management/latest/APIReference/API_Operations_AWS_Price_List_Service.html

https://aws.amazon.com/blogs/aws/aws-price-list-api-update-regional-price-lists/



## example

show rate (USD_JPY)
```console
$ awspricing ex
113.941498
```
using http://free.currencyconverterapi.com/api/v5/convert?q=USD_JPY&compact=y  
cache 2days


```console
$ awspricing ec2 c5.xlarge
CPU: 3.0 Ghz MEM: 8 GiB NETWORK: Up to 10 Gigabit ecu: 17 vcpu: 4 processor: Intel Xeon Platinum 8124M
PRICE: OnDemand 0.2140000000 USD (24.102820 JP) / Hrs
DESCRIPTION: $0.214 per On Demand Linux c5.xlarge Instance Hour
```

```console
$ awspricing ec2 c4.large --region us-west-2
CPU: 2.9 GHz MEM: 3.75 GiB NETWORK: Moderate ecu: 8 vcpu: 2 processor: Intel Xeon E5-2666 v3 (Haswell)
PRICE: OnDemand 0.1000000000 USD (11.263000 JP) / Hrs
DESCRIPTION: $0.1 per On Demand Linux c4.large Instance Hour
```
```console
$ awspricing ec2 c4.large --region ap-northeast-1
CPU: 2.9 GHz MEM: 3.75 GiB NETWORK: Moderate ecu: 8 vcpu: 2 processor: Intel Xeon E5-2666 v3 (Haswell)
PRICE: OnDemand 0.1260000000 USD (14.191380 JP) / Hrs
DESCRIPTION: $0.126 per On Demand Linux c4.large Instance Hour
```

```console
$ awspricing rds db.r4.xlarge
CPU: 2.3 GHz MEM: 30.5 GiB NETWORK: Up to 10 Gigabit vcpu: 4 processor: Intel Xeon E5-2686 v4 (Broadwell)
PRICE: OnDemand 0.5700000000 USD (64.199100 JP) / Hrs
DESCRIPTION: $0.570 per RDS db.r4.xlarge Single-AZ instance hour (or partial hour) running MySQL
```

```console
$ awspricing rds r4.xlarge
CPU: 2.3 GHz MEM: 30.5 GiB NETWORK: Up to 10 Gigabit vcpu: 4 processor: Intel Xeon E5-2686 v4 (Broadwell)
PRICE: OnDemand 0.5700000000 USD (64.199100 JP) / Hrs
DESCRIPTION: $0.570 per RDS db.r4.xlarge Single-AZ instance hour (or partial hour) running MySQL
```
