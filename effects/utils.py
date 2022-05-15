import pandas as pd
import numpy as np

def gaussian_mask(x, y, shape, amp=1, sigma=15):
    """
    Returns an array of shape, with values based on
    amp * exp(-((i-x)**2 +(j-y)**2) / (2 * sigma ** 2))
    :param x: float
    :param y: float
    :param shape: tuple
    :param amp: float
    :param sigma: float
    :return: array
    """
    xv, yv = np.meshgrid(np.arange(shape[1]), np.arange(shape[0]))
    g = amp * np.exp(-((xv - x) ** 2 + (yv - y) ** 2) / (2 * sigma ** 2))
    return g


def default(value, default_value):
    """
    Returns default_value if value is None, value otherwise
    """
    if value is None:
        return default_value
    return value