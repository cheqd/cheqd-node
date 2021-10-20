import string
import random


def get_alphanumeric_string(length):
    return ''.join(random.choice(string.ascii_letters + string.digits) for _ in range(length))
