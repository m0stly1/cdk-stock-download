import json
import datetime
import yfinance as yf
import boto3

def handler(event, context):
    end_date = datetime.datetime.now()
    start_date = end_date - datetime.timedelta(days=200)

    stocks = ["AAPL", "NFLX", "SPOT", "DISK", "TSLA"]

    for stock in stocks:
        data = yf.download(stock, start=start_date, end=end_date)
        data.to_csv(f"/tmp/{stock}.csv")
        s3 = boto3.resource("s3")
        s3.Bucket("awesome-stock-data-bucket").upload_file(f"/tmp/{stock}.csv", f"{stock}.csv")
