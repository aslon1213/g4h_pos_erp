import asyncio
import aiohttp
import time
import os

# def timeit(func):
#     async def wrapper(*args, **kwargs):
#         start_time = time.time()
#         result = await func(*args, **kwargs)
#         end_time = time.time()
#         print(f"Function {func.__name__} took {end_time - start_time:.4f} seconds")
#         return result

#     return wrapper


# @timeit
async def fetch_data(session: aiohttp.ClientSession, url, token):
    headers = {
        "Authorization": token,
        "Content-Type": "application/json",
    }
    start_time = time.time()
    async with aiohttp.request("GET", url, headers=headers) as response:
        result = await response.content.read()
        print(result)
    # async with session.request(url, headers=headers) as response:
    #     result = await response.content.read()
    #     print(result)

    return time.time() - start_time


async def main():
    token = os.getenv("TOKEN")
    host = os.getenv("HOST")
    if token == "" or host == "":
        print("TOKEN or HOST is not set")
        return

    async with aiohttp.ClientSession() as session:
        for _ in range(3):  # Run 10 batches
            tasks = []
            for _ in range(1000):  # 1000 requests per batch
                url = f"{host}/api/journals/649e78a656b78aefd50372e4"
                tasks.append(fetch_data(session, url, token))
            results = await asyncio.gather(*tasks)
            # calculate the average time
            average_time = sum(results) / len(results)
            print(f"Average time: {average_time:.4f} seconds")
            # print(results)
            await asyncio.sleep(1)  # Wait for 1 second before the next batch


if __name__ == "__main__":
    asyncio.run(main())
