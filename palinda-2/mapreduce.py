import multiprocessing

from collections import defaultdict

import re

def map_fn(document):
    kv_pairs = []

    words = document.lower().split()

    for word in words:

        clean_word = re.sub(r'[^a-z0-9]', "", word)

        if clean_word:
            kv_pairs.append((clean_word, 1))
    
    return kv_pairs


def reduce(kv_pairs):
    key, values = kv_pairs

    total_count = sum(values)

    return (key, total_count)


def run_mapreduce():

    documents = ["The quick brown fox jumps over the lazy dog", "To dog barked at the fox", "We encourage you to discuss with your course friends, but do not share answers! Similarily, use of AI services 🤖 are great to help explain things, but please do not submit AI-generated solutions - you must be both responsible for your own solutions and be able to explain them under examination. If in doubt, refer to the AI Policy on Canvas",   "We encourage you to discuss with your course friends, but do not share answers! Similarily, use of AI services 🤖 are great to help explain things, but please do not submit AI-generated solutions - you must be both responsible for your own solutions and be able to explain them under examination. If in doubt, refer to the AI Policy on Canvas", "We encourage you to discuss with your course friends, but do not share answers! Similarily, use of AI services 🤖 are great to help explain things, but please do not submit AI-generated solutions - you must be both responsible for your own solutions and be able to explain them under examination. If in doubt, refer to the AI Policy on Canvas", "We encourage you to discuss with your course friends, but do not share answers! Similarily, use of AI services 🤖 are great to help explain things, but please do not submit AI-generated solutions - you must be both responsible for your own solutions and be able to explain them under examination. If in doubt, refer to the AI Policy on Canvas"  ]


    print("MAP")
    map_res = None

    with multiprocessing.Pool(processes=multiprocessing.cpu_count()) as pool:
        map_res = pool.map(map_fn, documents)

    print("shuffle & sort")

    grouped_data = defaultdict(list)

    for process_result in map_res:
        for key, value in process_result:
            grouped_data[key].append(value)
    
    shuffled = list(grouped_data.items())

    print("Reduce")

    final_res = None

    with multiprocessing.Pool(processes=multiprocessing.cpu_count()) as pool:
        final_res = pool.map(reduce, shuffled)


    for word, count in final_res:

        print(f"{word}: {count}")
    
if __name__ == "__main__":
    run_mapreduce()
