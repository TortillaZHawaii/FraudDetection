# Use an official Python runtime as a parent image
FROM python:3.11-slim-buster AS base

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container at /app
COPY . /app

# Install any needed packages specified in requirements.txt
RUN pip install --trusted-host pypi.python.org -r requirements.txt

FROM base AS final

# Define environment variables
ENV KAFKA_BOOTSTRAP_SERVERS kafka:9092
ENV KAFKA_TOPIC alerts

# Run main.py when the container launches
CMD ["python", "main.py"]