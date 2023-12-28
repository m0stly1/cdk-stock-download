# Start from the AWS Lambda base image for Python 3.10
FROM public.ecr.aws/lambda/python:3.10

# Install Python dependencies
COPY requirements.txt .
RUN python -m pip install -r requirements.txt --target /asset

# Copy the Lambda function code
COPY main.py /asset/

# Set the CMD to your handler
CMD ["main.handler"]
