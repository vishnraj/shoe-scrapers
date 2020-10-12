import scrapy, argparse, urllib, json
from scrapy.crawler import CrawlerProcess

class ShoeSpider(scrapy.Spider):
  shoe_type = ''
  filename = ''
  base_urls = {
    'nike': 'https://www.nike.com/w'
  }

  def __init__(self, shoe_type, filename):
    if shoe_type == '':
      print('Requires a non-empty type')
      exit(1)
    self.shoe_type = shoe_type
    self.filename = filename

  def start_requests(self):
    request = scrapy.Request(f'{self.base_urls["nike"]}?q={urllib.parse.quote(self.shoe_type)}', callback=self.parse)
    return [request]

  def parse(self, response):
    out = []
    for product in response.css('div.product-card__body'):
      data = {
        'name': product.css('div.product-card__title::text').get(),
        'price': product.css('div.product-price::text').get()
      }
      out.append(data)
      yield data
    
    with open(self.filename, 'w+', encoding='utf8') as json_file:
      json.dump(out, json_file, ensure_ascii=False)



if __name__ == '__main__':
  parser = argparse.ArgumentParser(description='Scrapes shoe sites for a given shoe type')
  parser.add_argument('-s', '--shoe_type', help='Shoe type to query for', required=True)
  parser.add_argument('-o', '--out_file', help='JSON File to write data to', required=True)
  args = parser.parse_args()

  print('Gather data for shoe type: %s' % args.shoe_type)
  
  process = CrawlerProcess()
  process.crawl(ShoeSpider, args.shoe_type, args.out_file)
  process.start()