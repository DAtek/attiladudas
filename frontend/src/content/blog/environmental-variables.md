---
title: "Using environmental variables in python with direnv and datek-app-utils"
description: "A library for reading environmental variables in python"
pubDate: "2024-08-27"
---

## Foreword

This is an opinionated article. When I say something is _better_ or _is not correct_, I mean it only in my opinion, not scientifically.

# The current recommended way of reading env variables in python

I'm sure you already had some business with environmental variables, especially if you're following the [_III. principle_](https://12factor.net/config) of a _12factor application_.

<br>

I'm also sure you've read a lot of tutorials about how to deal with env vars, and unfortunately most of them recommend using the [`python-dotenv`](https://pypi.org/project/python-dotenv/) library.

Just a couple examples:

- https://www.freecodecamp.org/news/python-env-vars-how-to-get-an-environment-variable-in-python/
- https://datagy.io/python-environment-variables/
- https://www.datacamp.com/tutorial/python-environment-variables?dc_referrer=https%3A%2F%2Fduckduckgo.com%2F
- https://www.geeksforgeeks.org/access-environment-variable-values-in-python/

<br>

This article is offering a better, correct way of reading environmental variables.

# Why using `python-dotenv` isn't correct according to the _12factor_ ?

According to the the [_III. principle_](https://12factor.net/config),

<br>

"_The twelve-factor app **stores config in environment variables** (often shortened to env vars or env). Env vars are easy to change between deploys without changing any code; **unlike config files**, there is little chance of them being checked into the code repo accidentally; and unlike custom config files, or other config mechanisms such as Java System Properties, they are a language- and OS-agnostic standard."_

<br>

Let's talk about the following example:

<br>

```python
from os import getenv
from dotenv import load_dotenv


def main():
    load_dotenv()  # <--- This line reads a file!
    database_url = getenv("APP_DATABASE_URL")
    ...
```

<br>

As I pointed out, `load_dotenv()` line is actually reading a file! So the config originally comes from a file, not from the environment! With this effort you could read the config directly from `.yml` or `.json` or `.toml` files, right?

<br>

The whole point of _12factor_ is that the application should not read the config from files, which means when the application starts, the config should already be found in the environment!

# The correct way

Thankfully, there is a great tool to save us: [`direnv`](https://direnv.net/) !

<br>

With `direnv` you can create a `.envrc` file, which will be run when you're entering the working directory in your terminal.

<br>

A very minimalistic, one-liner `.envrc` file:

```
dotenv
```

<br>

With this one line, the `direnv` hook will read your `.env` file and it will automatically export all environmental variables into your shell.

<br>

With this, your app won't break the _III. principle_ anymore:

```python
from os import getenv


def main():
    database_url = getenv("APP_DATABASE_URL")
    ...
```

## Using `datek-app-utils`

With [`datek-app-utils`](https://pypi.org/project/datek_app_utils/) you can:

- Define a struct for your env vars
- Validate the environment (all env vars are defined and have the correct type)
- Access the env vars via class properties

<br>

```python
from datek_app_utils.env_config.base import BaseConfig, validate_config

class Config(BaseConfig):
    APP_DATABASE_URL: str


def main():
    validate_config(Config)

    database_url = Config.APP_DATABASE_URL
    ...
```

<br>

I hope you liked this post and I hope the tools I showed here will help you writing better software and improve your developer experience.

# Special thanks to

- [Anthony Jackson](https://github.com/expelledboy) for showing me `direnv` back in 2021
- [Dhia Kennouche](https://github.com/kendhia) for reviewing this post
