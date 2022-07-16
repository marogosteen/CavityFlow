import glob
import io
import os

import numpy as np
import matplotlib.pyplot as plt

SKIPROW = 1
LOGDIR = "../log/"
OUTPUTDIR = "../image/"
HEIGHT = 252
WIDTH = 252
U_COL = 0
V_COL = 1
P_COL = 2


def read_data(path: str) -> list:
    cavity_data = []
    with open(path) as f:
        skip_row(f)
        s = f.readline()
        while s:
            line = list(map(float, s.strip().split(",")))
            cavity_data.append(line[2:])
            s = f.readline()
    return cavity_data


def skip_row(f: io.TextIOWrapper):
    for _ in range(SKIPROW):
        f.readline()


def reshape_2d(array: np.ndarray, omit: int = 1) -> np.ndarray:
    if omit < 1:
        raise ValueError()

    array = array.reshape(HEIGHT, WIDTH)
    array = array[:HEIGHT]
    for _ in range(HEIGHT % omit):
        array = (array[:-1, :] + array[1:, :]) / 2
    for _ in range(WIDTH % omit):
        array = (array[:, :-1] + array[:, 1:]) / 2

    new_array = array[::omit, ::omit]
    for skip in range(1, omit):
        new_array += array[skip::omit, skip::omit]
    new_array /= omit

    return new_array


omit = 6

if not os.path.exists(OUTPUTDIR):
    os.mkdir(OUTPUTDIR)
else:
    for fp in glob.glob(OUTPUTDIR+"*.jpg"):
        os.remove(fp)

files = glob.glob(LOGDIR+"*.csv")
for fp in files[:1]:
    cavity_array = np.array(read_data(fp))
    u = reshape_2d(cavity_array[:, U_COL], omit=omit)
    v = reshape_2d(cavity_array[:, V_COL], omit=omit)
    p = reshape_2d(cavity_array[:, P_COL], omit=omit)

    row, col = u.shape
    u /= (col/2)
    v /= (row/2)
    x, y = np.meshgrid(
        np.linspace(0, 0.02, row),
        np.linspace(0, 0.02, col)
    )

    magnitude = np.hypot(u, v)

    fig, ax = plt.subplots(figsize=(14, 10))
    pcf = ax.contourf(x, y, p, alpha=0.3, cmap="jet", vmax=0.65, vmin=0.25)
    # TODO 戻す
    qv = ax.quiver(
        x, y, u, v, magnitude, cmap="jet",
        scale_units="xy", scale=1
    )
    fname = os.path.basename(fp)[:-len(".jpg")]

    savepath = OUTPUTDIR + f"{fname}.jpg"
    plt.colorbar(pcf)
    plt.subplots_adjust()
    plt.savefig(savepath)
    print(savepath)
    plt.close()

print("draw done")
