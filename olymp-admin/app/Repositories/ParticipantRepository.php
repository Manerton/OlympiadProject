<?php

namespace App\Repositories;

use App\Components\ApiHelper;
use App\Services\ApiService;
use Illuminate\Support\Facades\Cookie;

class ParticipantRepository
{
    private ApiService $apiService;
    public function __construct(ApiService $apiService){
        $this->apiService = $apiService;
    }
    public function getByApiAll($page = 1, $limit = 10, $token = null)
    {
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $response = $this->apiService->get(
            ApiHelper::PARTICIPANT_URL_API,
            [
                'page' => $page,
                'limit' => $limit
            ],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'] ? $response['data']['data'] : [];
    }
    public function getByApiId($id , $token = null)
    {
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $response =  $this->apiService->get(
            ApiHelper::PARTICIPANT_URL_API . '/' . $id,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'];
    }
    public function getByApiUserId($id, $token = null)
    {
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $response = $this->apiService->get(
            ApiHelper::PARTICIPANT_BY_USER_URL_API . '/' . $id,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'];
    }
    public function getCount($token = null)
    {
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $response = $this->apiService->get(
            ApiHelper::PARTICIPANT_COUNT_URL_API,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'];
    }
}
