"""
Convert txt-formatted question pools to python or JSON.
"""

import json
import re
import sys


def one(*args):
    """Like any/none, but for exactly one."""
    found = False
    for itm in args:
        if itm and found:
            return False
        found = found or bool(itm)
    return found


class PoolParse(object):

    SUB_RE = re.compile(
        r'^SUBELEMENT (?P<number>[TGE]\d+) . (?P<description>.*?)$',
        flags=re.MULTILINE)
    TOPIC_RE = re.compile(
        r'^(?P<number>[TGE]\d+[A-Z]+) - (?P<description>.*?)$',
        flags=re.MULTILINE)
    # If this regex looks strange, it's due to question pool inconsistencies
    QUESTION_RE = re.compile(
        r'^(?P<topic>[TGE]\d+[A-Z]+)(?P<number>\d+) '
        r'\((?P<correct>[A-D])\)'
        r'( \[(?P<citation>.*?)\]?)? ?\n'
        r'(?P<question>.*?)\n'
        r'A\. (?P<a>.*?)\n'
        r'B\. (?P<b>.*?)\n'
        r'C\. (?P<c>.*?)\n'
        r'D\. (?P<d>.*?)\n'
        r'~',
        flags=re.MULTILINE)

    CP1252 = [
            ('\x96', '-'),
            ('\xa0', ' '),
            ('\x92', "'"),
            ('\x93', '"'),
            ('\x94', '"'),
            ]

    def __init__(self, text=None, file_obj=None, filename=None, verbose=False):
        if not one(text, file_obj, filename):
            raise ValueError('Need one of text, file_obj, or filename.')
        if filename:
            file_obj = open(filename, 'rb')
        if file_obj:
            text = file_obj.read()
        # Windows 1252, manual replacement
        self.text = text.replace('\r\n', '\n')
        for s, r in self.CP1252:
            self.text = self.text.replace(s, r)
        self.verbose = verbose

    def parse(self):
        self.find_subelements()

    def find_subelements(self):
        results = {}
        for item in self.SUB_RE.finditer(self.text):
            results[item.group('number')] = item.group('description')
        if self.verbose:
            print 'Parsed {} subelements'.format(len(results))
        return results

    def find_topics(self):
        results = {}
        for item in self.TOPIC_RE.finditer(self.text):
            results[item.group('number')] = item.group('description')
        if self.verbose:
            print 'Parsed {} topics'.format(len(results))
        return results

    def find_questions(self):
        results = {}
        for item in self.QUESTION_RE.finditer(self.text):
            key = item.group('topic') + item.group('number')
            qdata = {
                'topic': item.group('topic'),
                'number': item.group('number'),
                'correct': item.group('correct'),
                'citation': item.group('citation'),
                'question': item.group('question'),
                'A': item.group('a'),
                'B': item.group('b'),
                'C': item.group('c'),
                'D': item.group('d'),
            }
            results[key] = qdata
        if self.verbose:
            print 'Parsed {} questions'.format(len(results))
        return results


def find_missing_qs(pool):
    known = {
            'T0A': 11,
            'T0B': 12,
            'T0C': 13,
            'T1A': 14,
            'T1B': 13,
            'T1C': 14,
            'T1D': 12,
            'T1E': 12,
            'T1F': 13,
            'T2A': 12,
            'T2B': 13,
            'T2C': 12,
            'T3A': 11,
            'T4A': 12,
            'T4B': 12,
            'T5A': 12,
            'T5B': 13,
            'T5C': 13,
            'T5D': 12,
            'T6A': 11,
            'T6B': 12,
            'T6C': 13,
            'T6D': 12,
            'T7A': 11,
            'T7B': 12,
            'T7C': 13,
            'T7D': 12,
            'T8A': 11,
            'T8B': 11,
            'T8C': 13,
            'T8D': 11,
            'T9A': 14,
            'T9B': 11,
            }
    all_questions = set()
    for topic, num in known.iteritems():
        for q in xrange(1, num+1):
            all_questions.add('%s%02d' % (topic, q))

    return all_questions - set(pool.find_questions().keys())


if __name__ == "__main__":
    parser = PoolParse(filename=sys.argv[1])
