from diplom_parser import logger

import constants
import requests


def get_olympiads() -> dict[str: int]:
    response = requests.get(f'http://{constants.API_HOST}:{constants.API_PORT}/api/v1/olympiads')
    if response.status_code != 200:
        logger.error(f'Failed to fetch olympiads. Status: {response.status_code}, Response: {response.text.strip()}')
        return dict()

    response.encoding = 'utf-8'

    olympiads = {}
    for olympiad in response.json():
        olympiads.setdefault(olympiad['name'].lower(), []).append(
            {
                'olympiad_id': olympiad['olympiad_id'],
                'profile': olympiad['profile'].lower()
            }
        )

    return olympiads