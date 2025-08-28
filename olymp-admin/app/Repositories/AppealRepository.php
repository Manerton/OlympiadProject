<?php

namespace App\Repositories;

use App\Components\Dictionaries\StatusDictionary;
use App\Models\Appeal;

class AppealRepository
{
    public function get($id)
    {
        return Appeal::find($id);
    }
    public function getAll()
    {
        return Appeal::all();
    }
    public function getByTaskId($taskId){
        return Appeal::where('task_id', $taskId)->get();
    }
    public function getByUserId($userId){
        return Appeal::where('user_id', $userId)->get();
    }
    public function delete($id)
    {
        return Appeal::destroy($id);
    }
    public function create($userId, $taskId, $reason){
        return Appeal::create([
            'user_id' => $userId,
            'task_id' => $taskId,
            'reason' => $reason,
            'status' => StatusDictionary::AWAITING
        ]);
    }
    public function changeStatus(Appeal $appeal, $status){
        $appeal->status = $status;
        $this->save($appeal);
    }
    public function save(Appeal $appeal){
        return $appeal->save();
    }
}
