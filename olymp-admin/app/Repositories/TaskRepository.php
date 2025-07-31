<?php

namespace App\Repositories;

use App\Models\Task;

class TaskRepository
{
    public function get($id)
    {
        return Task::where(['id' => $id])->first();
    }
    public function getAll(){
        return Task::all();
    }
    public function save(Task $task){
        return $task->save();
    }
    public function delete(Task $task){
        return $task->delete();
    }
    public function getByEventId($eventId){
        return Task::where(['event_id' => $eventId])->get();
    }
    public function create($eventId, $number, $points)
    {
        $task = new Task();
        $task->event_id = $eventId;
        $task->number = $number;
        $task->max_points = $points;
        $this->save($task);
        return $task->id;
    }
}
