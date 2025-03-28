FROM php:7.4-fpm


# 安装必要的系统依赖和 PHP 扩展
RUN apt-get update && apt-get install -y \
  git \
  gcc \
  make \
  re2c \
  libpcre3-dev \
  curl \
  libpng-dev \
  libjpeg-dev \
  libzip-dev \
  libgmp-dev \
  git \
  unzip \
  libcurl4-openssl-dev \
  libfreetype6-dev \
  build-essential \
  re2c \
  libjson-c-dev \
  && apt-get clean \
  && docker-php-ext-install zip bcmath curl sockets \
  && docker-php-ext-configure gd --with-freetype --with-jpeg \
  && docker-php-ext-install gd pdo mysqli pdo_mysql \
  && curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer \
  && pecl install redis && docker-php-ext-enable redis \
  && pecl install yaconf && docker-php-ext-enable yaconf \
  && echo "[yaconf]\nextension=yaconf.so\nyaconf.directory=/var/www/html/config" > /usr/local/etc/php/conf.d/yaconf.ini \
  # 安装 Zephir
  && composer global require phalcon/zephir \
  && ln -s /root/.composer/vendor/bin/zephir /usr/local/bin/zephir \
  # 将 Composer 的 bin 目录添加到 PATH
  && echo 'export PATH="$HOME/.composer/vendor/bin:$PATH"' >> ~/.bashrc \
  # 安装 php-zephir-parser（Zephir 依赖的扩展）
  && git clone https://github.com/zephir-lang/php-zephir-parser.git /tmp/php-zephir-parser \
  && cd /tmp/php-zephir-parser && phpize && ./configure && make && make install \
  && echo "extension=zephir_parser.so" > /usr/local/etc/php/conf.d/zephir_parser.ini \
  && rm -rf /tmp/php-zephir-parser

# 设置环境变量，确保 PATH 在容器启动时生效
ENV PATH="/root/.composer/vendor/bin:${PATH}"

# 设置工作目录
WORKDIR /var/www/html