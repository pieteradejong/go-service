from pyspark.sql import SparkSession
from pyspark.sql.functions import expr

# import os
# os.environ['PYSPARK_SUBMIT_ARGS'] = '--packages org.apache.spark:spark-streaming-kafka-0-10_2.12:3.5.1'

spark = SparkSession.builder.appName("EmojiCount").getOrCreate()

df = spark \
    .readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "kafka:9092") \
    .option("subscribe", "reaction-emoji-submission") \
    .load()

emojiCounts = df.groupBy("value").count()

query = emojiCounts \
    .selectExpr("CAST(value AS STRING) AS key", "CAST(count AS STRING) AS value") \
    .writeStream \
    .outputMode("complete") \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "kafka:9092") \
    .option("topic", "reaction-emoji-counts") \
    .option("checkpointLocation", "/opt/spark-apps/checkpoints") \
    .start()

query.awaitTermination()
