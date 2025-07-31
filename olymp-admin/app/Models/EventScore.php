<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;

class EventScore extends Model
{
    use HasFactory;
    public $eventAPI;
    protected $table = 'event_score';
    protected $fillable = [
        'event_id',
        'winner_score',
        'prize_score'
    ];
    public function setScore($winnerScore, $prizeScore){
        $this->winner_score = $winnerScore;
        $this->prize_score = $prizeScore;
    }
}
