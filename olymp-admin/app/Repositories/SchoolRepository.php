<?php

namespace App\Repositories;

use App\Components\ApiHelper;
use App\Services\ApiService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cookie;

class SchoolRepository
{
    private ApiService $apiService;
    public function __construct(
        ApiService $apiService
    )
    {
        $this->apiService = $apiService;
    }

    public function getByApiAll($page = 1, $limit = 10)
    {
        $token = null;
        if (isset(Cookie::getQueuedCookies()[0])) {
            $queuedCookie = Cookie::getQueuedCookies()[0];
            $token = json_decode($queuedCookie->getValue(), true)['token'] ?? null;
        }
        if (is_null($token) && Cookie::get('username')) {
            $token = json_decode(Cookie::get('username'), true)['token'] ?? null;
        }
        $response = $this->apiService->get(
            ApiHelper::SCHOOL_URL_API,
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
    public function getByApiId($id)
    {
        $token = null;
        if (isset(Cookie::getQueuedCookies()[0])) {
            $queuedCookie = Cookie::getQueuedCookies()[0];
            $token = json_decode($queuedCookie->getValue(), true)['token'] ?? null;
        }
        if (is_null($token) && Cookie::get('username')) {
            $token = json_decode(Cookie::get('username'), true)['token'] ?? null;
        }
        $response = $this->apiService->get(
            ApiHelper::SCHOOL_URL_API . '/' . $id,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'];
    }
    public function getCount()
    {
        $token = null;
        if (isset(Cookie::getQueuedCookies()[0])) {
            $queuedCookie = Cookie::getQueuedCookies()[0];
            $token = json_decode($queuedCookie->getValue(), true)['token'] ?? null;
        }
        if (is_null($token) && Cookie::get('username')) {
            $token = json_decode(Cookie::get('username'), true)['token'] ?? null;
        }
        $response = $this->apiService->get(
            ApiHelper::SCHOOL_COUNT_URL_API,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'];
    }
}
