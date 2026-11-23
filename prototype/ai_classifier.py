import time
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression

texts = ["Win a free prize click now urgent!", "Meeting notes from Thursday sync attached."] * 500
labels = [1, 0] * 500
vec = TfidfVectorizer().fit(texts)
clf = LogisticRegression().fit(vec.transform(texts), labels)

sample = ["Urgent: verify your invoice payment details immediately."] * 1000
t0 = time.perf_counter()
clf.predict(vec.transform(sample))
ms = (time.perf_counter() - t0) * 1000
print(f"Total: {ms:.2f} ms | Avg per message: {ms/len(sample):.3f} ms")

