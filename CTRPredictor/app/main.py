import falcon
import json
import pandas as pd
from sklearn.externals import joblib

class HelloResource:
    def on_get(self, req, resp):
        resp.status = falcon.HTTP_200
        resp.body = ('\nhello world\n\n')

class PredictResource:
    def on_post(self, req, resp):

        body = req.stream.read()
        data = json.loads(body.decode('utf-8'), cls=Decoder)
        advId = data['advertiserId']
        floorPrice = data['floorPrice']
        # user = data['user']
        # site = data['site']

        data = [floorPrice, advId] # data = [floorPrice, site, hash(user), advId]
        df = pd.DataFrame([data], columns=['floorPrice', 'advertiserId']) # columns=['floorPrice', 'site', 'user', 'advertiserId'])
        ctr = loaded_model.predict_proba(df)[:, 1]
        # print(ctr)

        resp.status = falcon.HTTP_200
        resp.body =  json.dumps({"ctr":str(ctr[0])}) # resp.body = str(ctr[0]).encode('utf-8','ignore')

class Decoder(json.JSONDecoder):
    def decode(self, s):
        result = super(Decoder, self).decode(s)
        return self._decode(result)
    def _decode(self, o):
        if isinstance(o, str):
            try:
                return int(o)
            except ValueError:
                try:
                    return float(o)
                except ValueError:
                    return o
        elif isinstance(o, dict):
           return {k: self._decode(v) for k, v in o.items()}
        elif isinstance(o, list):
            return [self._decode(v) for v in o]
        else:
            return o

app = falcon.API()
app.add_route('/hello', HelloResource())
app.add_route('/predict', PredictResource())

loaded_model = joblib.load('model_ad_fp.pkl')