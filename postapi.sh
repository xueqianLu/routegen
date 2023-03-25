#!/bin/bash
curl -H "Content-Type:application/json" -X POST --data '{"token0":"0xA4DCaA2f100c7Ed0ee1C5EbF0781334a9F21f05F","token1":"0xbb4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"}' http://127.0.0.1:9800/defiroute/api/v1/route
