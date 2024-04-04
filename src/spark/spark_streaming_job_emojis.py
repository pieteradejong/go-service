from pyspark.sql import SparkSession
from pyspark.sql.functions import from_json, window, col, struct, to_json
from pyspark.sql.types import StructType, StructField, StringType, LongType

schema = StructType([
    StructField("user_id", LongType(), True),
    StructField("timestamp", LongType(), True),
    StructField("emoji", StringType(), True)
])

spark = SparkSession.builder.appName("EmojiCount").getOrCreate()

df = spark \
    .readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "kafka:9092") \
    .option("subscribe", "reaction-emoji-submission") \
    .option("includeTimestamp", "true") \
    .load()

df = df.selectExpr("CAST(value AS STRING) as json_str") \
    .select(from_json("json_str", schema).alias("data")).select("data.*")

emojiCounts = df \
    .withColumn("timestamp", (col("timestamp") / 1000).cast("timestamp")) \
    .groupBy(
        window(col("timestamp"), "2 seconds"),
        col("emoji")
    ) \
    .count()

query = emojiCounts \
    .select(
        col("emoji").alias("key"),
        to_json(struct(col("emoji"), col("count"))).alias("value")
    ) \
    .writeStream \
    .outputMode("update") \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "kafka:9092") \
    .option("topic", "reaction-emoji-counts") \
    .option("checkpointLocation", "/opt/spark-apps/checkpoints") \
    .trigger(processingTime='2 seconds') \
    .start()

query.awaitTermination()
