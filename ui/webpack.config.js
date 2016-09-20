const path = require('path');
const webpack = require('webpack');
const merge = require('webpack-merge');
const HtmlWebpackPlugin = require('html-webpack-plugin');
const autoprefixer = require('autoprefixer');
const ExtractTextPlugin = require('extract-text-webpack-plugin');
const Copy = require('copy-webpack-plugin');

// detemine build env
const TARGET_ENV = process.env.npm_lifecycle_event === 'build' ? 'production' : 'development';

// common webpack config
const commonConfig = {

  output: {
    path: path.resolve(__dirname, 'dist/'),
    filename: '[hash].js',
  },

  resolve: {
    modulesDirectories: ['node_modules'],
    extensions: ['', '.js', '.elm'],
  },

  module: {
    noParse: /\.elm$/,
  },

  plugins: [
    new HtmlWebpackPlugin({
      template: 'src/index.html',
      inject: 'body',
      filename: 'index.html',
    }),
  ],

  postcss: [autoprefixer({
    browsers: ['last 2 versions'],
  })],

};

// additional webpack settings for local env (when invoked by 'npm start')
if (TARGET_ENV === 'development') {
  console.log('Serving locally...');

  module.exports = merge(commonConfig, {

    entry: [
      'webpack-dev-server/client?http://localhost:8000',
      path.join(__dirname, 'src/index.js'),
    ],

    devServer: {
      stats: 'errors-only',
      inline: true,
      progress: true,
      proxy: {
        '/api/*': 'http://localhost:8080',
      },
    },


    module: {
      loaders: [{
        test: /\.elm$/,
        exclude: [/elm-stuff/, /node_modules/],
        loader: 'elm-hot!elm-webpack',
      }, {
        test: /\.(css|scss)$/,
        loaders: [
          'style-loader',
          'css-loader',
          'postcss-loader',
          'sass-loader',
        ],
      }, {
        test: /\.json$/,
        loader: 'json-loader',
      }, {
        test: /\.jsx?$/,
        exclude: /(node_modules|bower_components)/,
        loader: 'babel', // 'babel-loader' is also a legal name to reference
        query: {
          presets: ['es2015'],
          cacheDirectory: true,
        },
      }],
    },

  });
}

// additional webpack settings for prod env (when invoked via 'npm run build')
if (TARGET_ENV === 'production') {
  console.log('Building for prod...');

  module.exports = merge(commonConfig, {
    entry: path.join(__dirname, 'src/index.js'),

    module: {
      loaders: [{
        test: /\.elm$/,
        exclude: [/elm-stuff/, /node_modules/],
        loader: 'elm-webpack',
      }, {
        test: /\.(css|scss)$/,
        loader: ExtractTextPlugin.extract('style-loader', [
          'css-loader',
          'postcss-loader',
          'sass-loader',
        ]),
      }, {
        test: /\.json$/,
        loader: 'json-loader',
      }, {
        test: /\.jsx?$/,
        exclude: /(node_modules|bower_components)/,
        loader: 'babel', // 'babel-loader' is also a legal name to reference
        query: {
          presets: ['es2015'],
          cacheDirectory: true,
        },
      }],
    },

    plugins: [
      new webpack.optimize.OccurenceOrderPlugin(),

      // extract CSS into a separate file
      new ExtractTextPlugin('./[hash].css', {
        allChunks: true,
      }),

      // new Copy([
      //   { from: 'src/images', to: 'images' },
      //   { from: 'src/favicon.ico' },
      //   { from: 'src/manifest.json' },
      // ]),

      // minify & mangle JS/CSS
      new webpack.optimize.UglifyJsPlugin({
        minimize: true,
        compressor: {
          warnings: false,
        },
        // mangle:  true
      }),
    ],

  });
}
