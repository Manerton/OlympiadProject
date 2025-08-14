<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class Log extends Model
{
    use HasFactory;
    public $id;
    public $ip;
    public $client;
    public $url;
    public $timestamps;
    public function fill($ip, $client, $url)
    {
        $this->ip = $ip;
        $this->client = $client;
        $this->url = $url;
        $this->timestamps = date('Y-m-d H:i:s');
    }
}
